package main

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"log"
	"os"
	"strconv"
	"strings"

	cli "github.com/virtual-kubelet/node-cli"
	"github.com/virtual-kubelet/node-cli/opts"
	"github.com/virtual-kubelet/node-cli/provider"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/service"
	vkProvider "github.com/atlarge-research/opendc-emulate-kubernetes/services/virtual_kubelet/provider"
	vkService "github.com/atlarge-research/opendc-emulate-kubernetes/services/virtual_kubelet/services"
)

var (
	buildVersion = "N/A"
	buildTime    = "N/A"
	k8sVersion   = "v1.15.2" // This should follow the version of k8s.io/kubernetes we are importing
	providerName = "changeme"
)

func main() {
	connectionInfo := service.NewConnectionInfo("localhost", 8080, true)

	joinApateCluster(connectionInfo)
	//startVK()
	//startGRPC()
}

func joinApateCluster(connectionInfo *service.ConnectionInfo) {
	c := vkService.GetJoinClusterClient(connectionInfo)
	defer func() {
		_ = c.Conn.Close()
	}()

	res, err := c.Client.JoinCluster(context.Background(), &empty.Empty{})

	if err != nil {
		log.Fatalf("Unable to join apate cluster: %v", err)
	}

	log.Printf("Joined apate cluster with token %s and uuid %s", res.KubernetesJoinToken, res.NodeUUID)
}

func startVK() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = cli.ContextWithCancelOnSignal(ctx)

	o := opts.New()
	o.Provider = providerName
	o.Version = strings.Join([]string{k8sVersion, providerName, buildVersion}, "-")

	node, err := cli.New(ctx,
		cli.WithBaseOpts(o),
		cli.WithCLIVersion(buildVersion, buildTime),
		cli.WithProvider(providerName, func(cfg provider.InitConfig) (provider.Provider, error) {
			return vkProvider.CreateProvider(), nil
		}),
	)

	if err != nil {
		log.Fatal(err)
	}

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}

func startGRPC() {
	// Connection settings
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal("Port not found in env")
	}

	connectionInfo := service.NewConnectionInfo("localhost", port, true)

	// Service
	server := service.NewGRPCServer(connectionInfo)
	vkService.RegisterScenarioService(server)
	server.Serve()
}
