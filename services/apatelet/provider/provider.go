package provider

import (
	"context"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/normalization"
	"github.com/sirupsen/logrus"
	cli "github.com/virtual-kubelet/node-cli"
	logruscli "github.com/virtual-kubelet/node-cli/logrus"
	"github.com/virtual-kubelet/node-cli/opts"
	"github.com/virtual-kubelet/node-cli/provider"
	"github.com/virtual-kubelet/virtual-kubelet/log"
	logruslogger "github.com/virtual-kubelet/virtual-kubelet/log/logrus"
	"strconv"
	"time"
)

var (
	k8sVersion = "v1.15.2" // This should follow the version of k8s.io/kubernetes we are importing
)

func CreateProvider(ctx context.Context, res *normalization.NodeResources, port int, port2 int) (*cli.Command, error) {
	logger := logrus.StandardLogger()

	log.L = logruslogger.FromLogrus(logrus.NewEntry(logger))
	logConfig := &logruscli.Config{LogLevel: "info"}
	op, err := opts.FromEnv()
	if err != nil {
		return nil, err
	}

	op.KubeConfigPath = "/tmp/apate/config"
	op.ListenPort = int32(port)
	op.MetricsAddr = ":" + strconv.Itoa(port2)
	op.Provider = "apatelet"
	op.PodSyncWorkers = 10
	op.InformerResyncPeriod = time.Second
	op.KubeNamespace = "kube-system"

	nodeInfo := cluster.NewNode("apatelet", "agent", "apatelet-"+res.UUID.String(), k8sVersion)

	node, err := cli.New(ctx,
		cli.WithProvider("apatelet", func(cfg provider.InitConfig) (provider.Provider, error) {
			cfg.DaemonPort = int32(port)
			return NewProvider(res, cfg, nodeInfo), nil
		}),
		cli.WithBaseOpts(op),
		cli.WithPersistentFlags(logConfig.FlagSet()),
		cli.WithPersistentPreRunCallback(func() error {
			return logruscli.Configure(logConfig, logger)
		}),
	)

	return node, err
}
