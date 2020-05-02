package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/services"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/store"
)

const (
	// ListenAddress is the address the control plane will listen on
	ListenAddress = "CP_LISTEN_ADDRESS"

	// ListenPort is the port the control plane will listen on
	ListenPort = "CP_LISTEN_PORT"

	// ManagedClusterConfig is the path to the config of the cluster manager, if applicable
	ManagedClusterConfig = "CP_K8S_CONFIG"

	// ExternalIP can be used to override the IP the control plane will give to apatelets to connect to
	ExternalIP = "CP_EXTERNAL_IP"
)

func init() {
	// Enable line numbers in logging
	// Enables date time flags & file name + line
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	log.Println("Starting Apate control plane")

	// Get external ip
	externalIP, err := getExternalAddress()

	if err != nil {
		log.Fatalf("Error while starting control plane: %s", err.Error())
	}

	log.Printf("External IP for control plane: %s", externalIP)

	// Create kubernetes cluster
	log.Println("Starting kubernetes control plane")
	managedKubernetesCluster := createCluster(getEnv(ManagedClusterConfig, "/tmp/apate/manager"))

	// Create apate cluster state
	createdStore := store.NewStore()

	// Start gRPC server
	server, err := createGRPC(&createdStore, managedKubernetesCluster.KubernetesCluster)

	if err != nil {
		log.Fatalf("Error while starting control plane: %s", err.Error())
	}

	log.Printf("Now accepting requests on %s:%d\n", server.Conn.Address, server.Conn.Port)

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

// TODO: Maybe use multiple detection methods later?
func getExternalAddress() (string, error) {
	// Check for external IP override
	if val, ok := os.LookupEnv(ExternalIP); ok {
		return val, nil
	}

	// Check for IP in interface addresses
	addresses, err := net.InterfaceAddrs()

	if err != nil {
		return "", err
	}

	for _, address := range addresses {
		if strings.Contains(address.String(), "172.17.") {
			ip := strings.Split(address.String(), "/")[0]

			return ip, nil
		}
	}

	// Default to localhost
	return "localhost", nil
}

func createGRPC(createdStore *store.Store, kubernetesCluster cluster.KubernetesCluster) (*service.GRPCServer, error) {
	// Retrieve from environment
	listenAddress := getEnv(ListenAddress, "0.0.0.0")
	listenPort, err := strconv.Atoi(getEnv(ListenPort, "8085"))

	if err != nil {
		return nil, err
	}

	// Connection settings
	connectionInfo := service.NewConnectionInfo(listenAddress, listenPort, false)

	// Create gRPC server
	server := service.NewGRPCServer(connectionInfo)

	// Add services
	services.RegisterStatusService(server, createdStore)
	services.RegisterScenarioService(server, createdStore)
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

func getEnv(key, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return def
}
