package main

import (
	"context"
	"log"
	"strings"

	cli "github.com/virtual-kubelet/node-cli"
	"github.com/virtual-kubelet/node-cli/opts"
	"github.com/virtual-kubelet/node-cli/provider"

	ourprovider "github.com/atlarge-research/opendc-emulate-kubernetes/services/virtual_kubelet/provider"
)

var (
	buildVersion = "N/A"
	buildTime    = "N/A"
	k8sVersion   = "v1.15.2" // This should follow the version of k8s.io/kubernetes we are importing
	providerName = "changeme"
)

func main() {
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
			return ourprovider.CreateProvider(), nil
		}),
	)

	if err != nil {
		log.Fatal(err)
	}

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
