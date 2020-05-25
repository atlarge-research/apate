package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/run"

	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/cluster/watchdog"

	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/crd/node"

	nodeconfigurationv1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/nodeconfiguration/v1"
	podconfigurationv1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/podconfiguration/v1"

	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/kubectl"

	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/network"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"

	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/service"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/services"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/store"
)

func init() {
	// Enable line numbers in logging
	// Enables date time flags & file name + line
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

func panicf(err error) {
	log.Panicf("An error occurred while starting the controlplane: %+v\n", err)
}

func main() {
	log.Println("starting Apate control plane")

	ctx := context.Background()

	// Get external connection information
	externalInformation, err := createExternalConnectionInformation()
	if err != nil {
		panicf(errors.Wrap(err, "failed to get external connection information"))
	}

	// Register runners
	registerRunners()

	// Create kubernetes cluster
	managedKubernetesCluster, err := createCluster(env.ControlPlaneEnv().ManagerConfigLocation)
	if err != nil {
		panicf(errors.Wrap(err, "failed to create cluster"))
	}

	// Create apate cluster state
	createdStore := store.NewStore()

	// Save the kubeconfig in the store
	if err = createdStore.SetKubeConfig(*managedKubernetesCluster.KubeConfig); err != nil {
		panicf(errors.Wrap(err, "failed to set Kubeconfig"))
	}

	// Create CRDs
	if err = createCRDs(managedKubernetesCluster); err != nil {
		panicf(errors.Wrap(err, "failed to create CRDs"))
	}

	// TODO: Remove later, seems to give k8s some breathing room for crd
	time.Sleep(time.Second)

	// Create node informer
	stopInformer := make(chan struct{})
	if err = node.CreateNodeInformer(ctx, managedKubernetesCluster.KubeConfig, &createdStore, externalInformation, stopInformer); err != nil {
		panicf(errors.Wrap(err, "failed to create node informer"))
	}

	// Create prometheus stack
	createPrometheus := env.ControlPlaneEnv().PrometheusStackEnabled
	if createPrometheus {
		go kubectl.CreatePrometheusStack(managedKubernetesCluster.KubeConfig)
	}

	// Start gRPC server
	server, err := createGRPC(&createdStore, managedKubernetesCluster.Cluster, externalInformation)
	if err != nil {
		panicf(errors.Wrap(err, "failed to start GRPC server"))
	}

	log.Printf("now accepting requests on %s:%d\n", server.Conn.Address, server.Conn.Port)

	kubeConfigLocation := env.ControlPlaneEnv().KubeConfigLocation
	if err = ioutil.WriteFile(kubeConfigLocation, managedKubernetesCluster.Cluster.KubeConfig.Bytes, 0600); err != nil {
		panicf(errors.Wrap(err, "failed to write Kubeconfig to tempfile"))
	}

	// Handle signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Start serving request
	go func() {
		err := server.Serve()

		if err != nil {
			panicf(errors.Wrap(err, "failed to start gRPC server"))
		}
	}()

	// Start watchdog
	watchdog.StartWatchDog(ctx, time.Second*30, &createdStore, &managedKubernetesCluster.Cluster)

	// Stop the server on signal
	<-stop
	stopInformer <- struct{}{}
	shutdown(&createdStore, &managedKubernetesCluster, server)
	log.Printf("apate control plane stopped")
}

func createCRDs(managedKubernetesCluster kubernetes.ManagedCluster) error {
	if err := podconfigurationv1.CreateInKubernetes(managedKubernetesCluster.KubeConfig); err != nil {
		return errors.Wrap(err, "failed to register pod CRD spec")
	}

	if err := nodeconfigurationv1.CreateInKubernetes(managedKubernetesCluster.KubeConfig); err != nil {
		return errors.Wrap(err, "failed to register node CRD spec")
	}

	return nil
}

func registerRunners() {
	var dockerRunner run.ApateletRunner = run.DockerRunner{}
	run.Registry.RegisterRunner(env.Docker, &dockerRunner)

	var routineRunner run.ApateletRunner = run.RoutineRunner{}
	run.Registry.RegisterRunner(env.Routine, &routineRunner)
}

func shutdown(store *store.Store, kubernetesCluster *kubernetes.ManagedCluster, server *service.GRPCServer) {
	log.Println("stopping Apate control plane")

	log.Println("stopping API")
	server.Server.Stop()

	// TODO: Actual cleanup for other nodes, for now just wipe state
	if err := (*store).ClearNodes(); err != nil {
		log.Printf("an error occurred while cleaning the apate store: %s", err.Error())
	}

	log.Println("stopping kubernetes control plane")
	if err := kubernetesCluster.Delete(); err != nil {
		log.Printf("an error occurred while deleting the kubernetes store: %s", err.Error())
	}
}

func getExternalAddress() (string, error) {
	// Check for external IP override
	override := env.ControlPlaneEnv().ExternalIP
	if override != env.ControlPlaneExternalIPDefault {
		return override, nil
	}

	res, err := network.GetExternalAddress()
	if err != nil {
		return "", errors.Wrap(err, "failed to retrieve external ip address from the 'network' module")
	}

	return res, nil
}

func createGRPC(createdStore *store.Store, kubernetesCluster kubernetes.Cluster, info *service.ConnectionInfo) (*service.GRPCServer, error) {
	// Retrieve from environment
	listenAddress := env.ControlPlaneEnv().ListenAddress

	// Connection settings
	connectionInfo := service.NewConnectionInfo(listenAddress, info.Port, false)

	// Create gRPC server
	server, err := service.NewGRPCServer(connectionInfo)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create new GRPC server with connection info %v", connectionInfo)
	}

	// Add services
	services.RegisterStatusService(server, createdStore)
	services.RegisterScenarioService(server, createdStore, info)
	services.RegisterClusterOperationService(server, createdStore, kubernetesCluster)
	services.RegisterHealthService(server, createdStore)

	return server, nil
}

func createCluster(managedClusterConfigPath string) (kubernetes.ManagedCluster, error) {
	log.Println("starting kubernetes control plane")

	cb := kubernetes.Default()
	c, err := cb.WithName("Apate").WithManagerConfig(managedClusterConfigPath).ForceCreate()
	if err != nil {
		return kubernetes.ManagedCluster{}, errors.Wrap(err, "failed to create new cluster")
	}

	numberOfPods, err := c.GetNumberOfPods("kube-system")
	if err != nil {
		return kubernetes.ManagedCluster{}, errors.Wrap(err, "failed to get number of pods from kubernetes")
	}

	log.Printf("There are %d pods in the cluster", numberOfPods)

	return c, nil
}

func createExternalConnectionInformation() (*service.ConnectionInfo, error) {
	// Get external ip
	externalIP, err := getExternalAddress()

	if err != nil {
		return nil, errors.Wrap(err, "failed to get external ip address")
	}

	// Get port
	listenPort := env.ControlPlaneEnv().ListenPort

	// Create external information
	log.Printf("external IP for control plane: %s", externalIP)

	return service.NewConnectionInfo(externalIP, listenPort, false), nil
}
