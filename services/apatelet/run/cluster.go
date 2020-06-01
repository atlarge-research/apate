package run

import (
	"context"
	"log"

	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/service"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/clients/controlplane"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes/kubeconfig"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"
	crdNode "github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/crd/node"
	crdPod "github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/crd/pod"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/scheduler"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
)

func joinApateCluster(ctx context.Context, connectionInfo *service.ConnectionInfo, listenPort int, kubeConfigPath string) (*kubeconfig.KubeConfig, *scenario.NodeResources, int64, error) {
	log.Println("Joining apate cluster")

	client, err := controlplane.GetClusterOperationClient(connectionInfo)
	if err != nil {
		return nil, nil, -1, errors.Wrap(err, "failed to get cluster operation client")
	}

	defer func() {
		closeErr := client.Conn.Close()
		if closeErr != nil {
			log.Printf("could not close connection: %v\n", closeErr)
		}
	}()

	cfg, res, startTime, err := client.JoinCluster(ctx, listenPort, kubeConfigPath)

	if err != nil {
		return nil, nil, -1, errors.Wrap(err, "failed to join cluster")
	}

	log.Printf("Joined apate cluster with resources: %v", res)

	return cfg, res, startTime, nil
}

func createInformers(config *kubeconfig.KubeConfig, st store.Store, stopInformer <-chan struct{}, sch *scheduler.Scheduler, res *scenario.NodeResources) error {
	err := crdPod.CreatePodInformer(config, &st, stopInformer, sch.WakeScheduler)
	if err != nil {
		return errors.Wrap(err, "failed creating crd pod informer")
	}

	err = crdNode.CreateNodeInformer(config, &st, res.Label, stopInformer, sch.WakeScheduler)
	if err != nil {
		return errors.Wrap(err, "failed creating crd node informer")
	}

	return nil
}
