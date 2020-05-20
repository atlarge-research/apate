package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/watchdog"

	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/crd/node"

	nodeconfigurationv1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/nodeconfiguration/v1"
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
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

func main() {
	log.Println("starting Apate control plane")

	ctx := context.Background()

	// Get external connection information
	externalInformation, err := createExternalConnectionInformation()

	if err != nil {
		log.Fatalf("error while starting control plane: %s", err.Error())
	}

	// Create kubernetes cluster
	log.Println("starting kubernetes control plane")
	managedKubernetesCluster := createCluster(env.RetrieveFromEnvironment(env.ManagedClusterConfig, env.ManagedClusterConfigDefault))

	// Create apate cluster state
	createdStore := store.NewStore()

	// Save the kubeconfig in the store
	if err = createdStore.SetKubeConfig(*managedKubernetesCluster.KubeConfig); err != nil {
		log.Fatal(err)
	}

	// Create CRDs
	if err = podconfigurationv1.CreateInKubernetes(managedKubernetesCluster.KubeConfig); err != nil {
		log.Fatal(err)
	}
	if err = nodeconfigurationv1.CreateInKubernetes(managedKubernetesCluster.KubeConfig); err != nil {
		log.Fatal(err)
	}

	// TODO: Remove later, seems to give k8s some breathing room for crd
	time.Sleep(time.Second)

	// Create node informer
	stopInformer := make(chan struct{})
	if err = node.CreateNodeInformer(ctx, managedKubernetesCluster.KubeConfig, &createdStore, externalInformation, stopInformer); err != nil {
		log.Fatal(err)
	}

	// Create prometheus stack
	createPrometheus := env.RetrieveFromEnvironment(env.PrometheusStackEnabled, env.PrometheusStackEnabledDefault)
	if strings.ToLower(createPrometheus) == "true" {
		go kubectl.CreatePrometheusStack(managedKubernetesCluster.KubeConfig)
	}

	// Start gRPC server
	server, err := createGRPC(&createdStore, managedKubernetesCluster.KubernetesCluster, externalInformation)
	if err != nil {
		log.Fatalf("error while starting control plane: %s", err.Error())
	}

	log.Printf("now accepting requests on %s:%d\n", server.Conn.Address, server.Conn.Port)

	if err = ioutil.WriteFile(os.TempDir()+"/apate/config", managedKubernetesCluster.KubernetesCluster.KubeConfig.Bytes, 0600); err != nil {
		log.Fatalf("error while starting control plane: %s", err.Error())
	}

	// Handle signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Start serving request
	go server.Serve()

	// Start watchdog
	watchdog.StartWatchDog(time.Second*30, &createdStore, &managedKubernetesCluster.KubernetesCluster)

	// Stop the server on signal
	<-stop
	stopInformer <- struct{}{}
	shutdown(&createdStore, &managedKubernetesCluster, server)
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

	return network.GetExternalAddress()
}

func createGRPC(createdStore *store.Store, kubernetesCluster cluster.KubernetesCluster, info *service.ConnectionInfo) (*service.GRPCServer, error) {
	// Retrieve from environment
	listenAddress := env.RetrieveFromEnvironment(env.ControlPlaneListenAddress, env.ControlPlaneListenAddressDefault)

	// Connection settings
	connectionInfo := service.NewConnectionInfo(listenAddress, info.Port, false)

	// Create gRPC server
	server, err := service.NewGRPCServer(connectionInfo)
	if err != nil {
		return nil, err
	}

	// Add services
	services.RegisterStatusService(server, createdStore)
	services.RegisterScenarioService(server, createdStore, info)
	services.RegisterClusterOperationService(server, createdStore, kubernetesCluster)
	services.RegisterHealthService(server, createdStore)

	return server, nil
}

func createCluster(managedClusterConfigPath string) cluster.ManagedCluster {
	cb := cluster.Default()
	c, err := cb.WithName("Apate").WithManagerConfig(managedClusterConfigPath).ForceCreate()
	if err != nil {
		log.Fatalf("an error occurred: %s", err.Error())
	}

	numberOfPods, err := c.GetNumberOfPods("kube-system")
	if err != nil {
		log.Fatalf("an error occurred: %s", err.Error())
	}

	log.Printf("There are %d pods in the cluster", numberOfPods)

	return c
}

func createExternalConnectionInformation() (*service.ConnectionInfo, error) {
	// Get external ip
	externalIP, err := getExternalAddress()

	if err != nil {
		return nil, err
	}

	// Get port
	listenPort, err := strconv.Atoi(env.RetrieveFromEnvironment(env.ControlPlaneListenPort, env.ControlPlaneListenPortDefault))

	if err != nil {
		return nil, err
	}

	// Create external information
	log.Printf("external IP for control plane: %s", externalIP)

	return service.NewConnectionInfo(externalIP, listenPort, false), nil
}
