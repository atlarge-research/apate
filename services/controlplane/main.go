package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/pkg/errors"

	podconfigurationv1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/podconfiguration/v1"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubectl"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/network"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/services"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/store"
)

func init() {
	// Enable line numbers in logging
	// Enables date time flags & file name + line
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func fatal(err error) {
	log.Fatalf("An error occurred while starting the controlplane: %+v\n", err)
}

func main() {
	log.Println("starting Apate control plane")

	// Get external connection information
	externalInformation, err := createExternalConnectionInformation()

	if err != nil {
		fatal(errors.Wrap(err, "failed to get external connection information"))
	}

	// Create kubernetes cluster
	log.Println("starting kubernetes control plane")
	managedKubernetesCluster, err := createCluster(env.RetrieveFromEnvironment(env.ManagedClusterConfig, env.ManagedClusterConfigDefault))
	if err != nil {
		fatal(errors.Wrap(err, "failed to create cluster"))
	}

	// Create apate cluster state
	createdStore := store.NewStore()

	// Save the kubeconfig in the store
	if err = createdStore.SetKubeConfig(*managedKubernetesCluster.KubeConfig); err != nil {
		fatal(errors.Wrap(err, "failed to set Kubeconfig"))
	}

	if err = podconfigurationv1.CreateInKubernetes(managedKubernetesCluster.KubeConfig); err != nil {
		fatal(errors.Wrap(err, "failed to register pod CRD spec"))
	}

	// Create prometheus stack
	createPrometheus := env.RetrieveFromEnvironment(env.PrometheusStackEnabled, env.PrometheusStackEnabledDefault)
	if strings.ToLower(createPrometheus) == "true" {
		go kubectl.CreatePrometheusStack(managedKubernetesCluster.KubeConfig)
	}

	// Start gRPC server
	server, err := createGRPC(&createdStore, managedKubernetesCluster.KubernetesCluster, externalInformation)
	if err != nil {
		fatal(errors.Wrap(err, "failed to start GRPC server"))
	}

	log.Printf("now accepting requests on %s:%d\n", server.Conn.Address, server.Conn.Port)

	if err = ioutil.WriteFile(os.TempDir()+"/apate/config", managedKubernetesCluster.KubernetesCluster.KubeConfig.Bytes, 0600); err != nil {
		fatal(errors.Wrap(err, "failed to write Kubeconfig to tempfile"))
	}

	// Handle signals
	signals := make(chan os.Signal, 1)
	stopped := make(chan bool, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		shutdown(&createdStore, &managedKubernetesCluster, server)
		stopped <- true
	}()

	// Start serving request
	go server.Serve()

	// Stop the server on signal
	<-stopped
	log.Printf("apate control plane stopped")
}

func shutdown(store *store.Store, kubernetesCluster *cluster.ManagedCluster, server *service.GRPCServer) {
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
	override := env.RetrieveFromEnvironment(env.ControlPlaneExternalIP, env.ControlPlaneExternalIPDefault)
	if override != env.ControlPlaneExternalIPDefault {
		return override, nil
	}

	res, err := network.GetExternalAddress()
	if err != nil {
		return "", errors.Wrap(err, "failed to retrieve external ip address from the 'network' module")
	}

	return res, nil
}

func createGRPC(createdStore *store.Store, kubernetesCluster cluster.KubernetesCluster, info *service.ConnectionInfo) (*service.GRPCServer, error) {
	// Retrieve from environment
	listenAddress := env.RetrieveFromEnvironment(env.ControlPlaneListenAddress, env.ControlPlaneListenAddressDefault)

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

func createCluster(managedClusterConfigPath string) (cluster.ManagedCluster, error) {
	cb := cluster.Default()
	c, err := cb.WithName("Apate").WithManagerConfig(managedClusterConfigPath).ForceCreate()
	if err != nil {
		return cluster.ManagedCluster{}, errors.Wrap(err, "failed to create new cluster")
	}

	numberOfPods, err := c.GetNumberOfPods("kube-system")
	if err != nil {
		return cluster.ManagedCluster{}, errors.Wrap(err, "failed to get number of pods from kubernetes")
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
	portstring := env.RetrieveFromEnvironment(env.ControlPlaneListenPort, env.ControlPlaneListenPortDefault)
	listenPort, err := strconv.Atoi(portstring)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to convert listening port to an integer (was %v)", portstring)
	}

	// Create external information
	log.Printf("external IP for control plane: %s", externalIP)

	return service.NewConnectionInfo(externalIP, listenPort, false), nil
}
