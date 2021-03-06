// Package run is the main package for the controlplane
package run

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/atlarge-research/apate/pkg/channel"

	"github.com/atlarge-research/apate/services/controlplane/crd/pod"

	"github.com/pkg/errors"

	"github.com/atlarge-research/apate/internal/kubectl"
	"github.com/atlarge-research/apate/internal/network"
	"github.com/atlarge-research/apate/internal/service"
	nodeconfigv1 "github.com/atlarge-research/apate/pkg/apis/nodeconfiguration/v1"
	podconfigv1 "github.com/atlarge-research/apate/pkg/apis/podconfiguration/v1"
	"github.com/atlarge-research/apate/pkg/env"
	"github.com/atlarge-research/apate/pkg/kubernetes"
	"github.com/atlarge-research/apate/pkg/runner"
	"github.com/atlarge-research/apate/services/controlplane/cluster/watchdog"
	"github.com/atlarge-research/apate/services/controlplane/crd/node"
	"github.com/atlarge-research/apate/services/controlplane/services"
	"github.com/atlarge-research/apate/services/controlplane/store"

	_ "k8s.io/client-go/plugin/pkg/client/auth" // Needed to connect with different providers, such as GKE
)

func init() {
	// Enable line numbers in logging
	// Enables date time flags & file name + line
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

func panicf(err error) {
	log.Panicf("An error occurred while starting the controlplane: %+v\n", err)
}

// StartControlPlane is the main control plane entrypoint
func StartControlPlane(ctx context.Context, registry *runner.Registry) {
	stop := make(chan os.Signal, 1)
	StartControlPlaneWithStopCh(ctx, registry, stop)
}

// StartControlPlaneWithStopCh starts the controlplane with a stop channel.
func StartControlPlaneWithStopCh(ctx context.Context, registry *runner.Registry, stopCh chan os.Signal) {
	cpEnv := env.ControlPlaneEnv()

	log.Println("starting Apate control plane")

	// Register runners
	registerRunners(registry)

	// Create kubernetes cluster
	cluster, err := createCluster()
	if err != nil {
		panicf(errors.Wrap(err, "failed to create cluster"))
	}

	var clusterAPI kubernetes.ClusterAPI = cluster

	// Create apate cluster state
	createdStore := store.NewStore()

	// Create CRDs
	if err = createCRDs(cluster); err != nil {
		panicf(errors.Wrap(err, "failed to create CRDs"))
	}

	// TODO: Remove later, seems to give k8s some breathing room for crd
	time.Sleep(time.Second)

	// Get external connection information
	externalInformation, err := createExternalConnectionInformation()
	if err != nil {
		panicf(errors.Wrap(err, "failed to get external connection information"))
	}

	// Create node informer
	stopInformer := channel.NewStopChannel()
	nodeHandler := node.NewHandler(&createdStore, registry, externalInformation, cluster)
	if err = node.WatchHandler(ctx, cluster.KubeConfig, nodeHandler, stopInformer.GetChannel()); err != nil {
		panicf(errors.Wrap(err, "failed to watch node handler"))
	}

	if err = pod.NoopWatchHandler(cluster.KubeConfig, stopInformer.GetChannel()); err != nil {
		panicf(errors.Wrap(err, "failed to watch pod handler"))
	}

	// Create prometheus stack
	createPrometheus := cpEnv.PrometheusEnabled
	if createPrometheus {
		go kubectl.CreatePrometheusStack(cluster.KubeConfig)
	}

	// Start gRPC server
	server, err := createGRPC(&createdStore, cluster, externalInformation, stopInformer)
	if err != nil {
		panicf(errors.Wrap(err, "failed to start GRPC server"))
	}
	externalInformation.Port = server.Conn.Port

	log.Printf("now accepting requests on %s:%d\n", server.Conn.Address, server.Conn.Port)

	kubeConfigLocation := cpEnv.KubeConfigLocation
	if err = ioutil.WriteFile(kubeConfigLocation, cluster.KubeConfig.Bytes, 0600); err != nil {
		panicf(errors.Wrap(err, "failed to write Kubeconfig to tempfile"))
	}

	// Handle signals
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)

	// Start serving request
	go func() {
		err := server.Serve()

		if err != nil {
			panicf(errors.Wrap(err, "failed to start gRPC server"))
		}
	}()

	// Start watchdog
	watchdog.StartWatchDog(ctx, time.Second*30, &createdStore, &clusterAPI)

	// Stop the server on signal
	select {
	case <-stopCh:
		//
	case <-ctx.Done():
		//
	}
	stopInformer.Close()
	shutdown(&createdStore, cluster, server)
	log.Printf("apate control plane stopped")
}

func createCRDs(cluster *kubernetes.Cluster) error {
	if err := podconfigv1.UpdateInKubernetes(cluster.KubeConfig, false); err != nil {
		return errors.Wrap(err, "failed to register pod CRD spec")
	}

	if err := nodeconfigv1.UpdateInKubernetes(cluster.KubeConfig, false); err != nil {
		return errors.Wrap(err, "failed to register node CRD spec")
	}

	return nil
}

func registerRunners(registry *runner.Registry) {
	var dockerRunner runner.ApateletRunner = &runner.DockerRunner{}
	registry.RegisterRunner(env.Docker, &dockerRunner)

	var routineRunner runner.ApateletRunner = &runner.RoutineRunner{}
	registry.RegisterRunner(env.Routine, &routineRunner)
}

func shutdown(store *store.Store, cluster *kubernetes.Cluster, server *service.GRPCServer) {
	log.Println("stopping Apate control plane")

	log.Println("stopping API")
	server.Server.Stop()

	// TODO: Actual cleanup for other nodes, for now just wipe state
	if err := (*store).ClearNodes(); err != nil {
		log.Printf("an error occurred while cleaning the apate store: %s", err.Error())
	}

	log.Println("stopping kubernetes control plane")
	if err := cluster.Shutdown(); err != nil {
		log.Printf("an error occurred while deleting the kubernetes cluster: %s", err.Error())
	}
}

func getExternalAddress() (string, error) {
	// Check for external IP override
	override := env.ControlPlaneEnv().ExternalIP
	if override != env.CPExternalIPDefault {
		return override, nil
	}

	res, err := network.GetExternalAddress()
	if err != nil {
		return "", errors.Wrap(err, "failed to retrieve external ip address from the 'network' module")
	}

	return res, nil
}

func createGRPC(createdStore *store.Store, kubernetesCluster *kubernetes.Cluster, info *service.ConnectionInfo, stopInformerCh *channel.StopChannel) (*service.GRPCServer, error) {
	// Retrieve from environment
	listenAddress := env.ControlPlaneEnv().ListenAddress

	// Connection settings
	connectionInfo := service.NewConnectionInfo(listenAddress, info.Port)

	// Create gRPC server
	server, err := service.NewGRPCServer(connectionInfo)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create new GRPC server with connection info %v", connectionInfo)
	}

	// Add services
	services.RegisterStatusService(server, createdStore)
	services.RegisterScenarioService(server, createdStore, info, stopInformerCh)
	services.RegisterClusterOperationService(server, createdStore, kubernetesCluster)
	services.RegisterHealthService(server, createdStore)

	return server, nil
}

func createCluster() (*kubernetes.Cluster, error) {
	log.Println("starting kubernetes control plane")

	cmh := kubernetes.NewClusterManagerHandler()

	c, err := cmh.NewCluster()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get cluster")
	}

	numberOfPods, err := c.GetNumberOfPods("kube-system")
	if err != nil {
		return nil, errors.Wrap(err, "failed to get number of pods from kubernetes")
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

	return service.NewConnectionInfo(externalIP, listenPort), nil
}
