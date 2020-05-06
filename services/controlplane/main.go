package main

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	ec "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"

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

func main() {
	log.Println("Starting Apate control plane")

	// Get external connection information
	externalInformation, err := createExternalConnectionInformation()

	if err != nil {
		log.Fatalf("Error while starting control plane: %s", err.Error())
	}

	// Create kubernetes cluster
	log.Println("Starting kubernetes control plane")
	managedKubernetesCluster := createCluster(env.RetrieveFromEnvironment(env.ManagedClusterConfig, env.ManagedClusterConfigDefault))

	// Create apate cluster state
	createdStore := store.NewStore()

	// Start gRPC server
	server, err := createGRPC(&createdStore, managedKubernetesCluster.KubernetesCluster, externalInformation)
	if err != nil {
		log.Fatalf("Error while starting control plane: %s", err.Error())
	}

	log.Printf("Now accepting requests on %s:%d\n", server.Conn.Address, server.Conn.Port)

	err = ioutil.WriteFile("/tmp/apate/config", managedKubernetesCluster.KubernetesCluster.KubeConfig, 0600)

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
	server.Serve()

	// Stop the server on signal
	<-stopped
	log.Printf("Apate control plane stopped")
}

func shutdown(store *store.Store, kubernetesCluster *cluster.ManagedCluster, server *service.GRPCServer) {
	log.Println("Stopping Apate control plane")

	log.Println("Stopping API")
	server.Server.Stop()

	// TODO: Actual cleanup for other nodes, for now just wipe state
	if err := (*store).ClearNodes(); err != nil {
		log.Printf("An error occurred while cleaning the apate store: %s", err.Error())
	}

	log.Println("Stopping kubernetes control plane")
	if err := kubernetesCluster.Delete(); err != nil {
		log.Printf("An error occurred while deleting the kubernetes store: %s", err.Error())
	}
}

// TODO: Maybe check for docker subnet first somehow, people can change it from 172.17.0.0/16 to something else after all..

// getExternalAddress will return the detected external IP address based on the env var, then network interfaces
// (it will look for the first 172.17.0.0/16 address), and finally a fallback on localhost
func getExternalAddress() (string, error) {
	// Check for external IP override
	override := env.RetrieveFromEnvironment(env.ControlPlaneExternalIP, env.ControlPlaneExternalIPDefault)
	if override != env.ControlPlaneExternalIPDefault {
		return override, nil
	}

	// Check for IP in interface addresses
	addresses, err := net.InterfaceAddrs()

	if err != nil {
		return "", err
	}

	// Get first 172.17.0.0/16 address, if any
	for _, address := range addresses {
		if strings.Contains(address.String(), ec.DockerAddressPrefix) {
			ip := strings.Split(address.String(), "/")[0]

			return ip, nil
		}
	}

	// Default to localhost
	return "localhost", nil
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
		log.Fatalf("An error occurred: %s", err.Error())
	}

	numberOfPods, err := c.GetNumberOfPods("kube-system")
	if err != nil {
		log.Fatalf("An error occurred: %s", err.Error())
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
	log.Printf("External IP for control plane: %s", externalIP)

	return service.NewConnectionInfo(externalIP, listenPort, false), nil
}
