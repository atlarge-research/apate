// Package run is the entry point of the actual apatelet
package run

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"

	"github.com/google/uuid"

	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/scheduler"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"

	"github.com/pkg/errors"

	crdNode "github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/crd/node"
	crdPod "github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/crd/pod"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes/kubeconfig"

	healthpb "github.com/atlarge-research/opendc-emulate-kubernetes/api/health"
	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/service"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/clients/controlplane"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/clients/health"
	vkProvider "github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/provider"
	vkService "github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/services"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
)

var once sync.Once

// StartApatelet starts the apatelet
func StartApatelet(apateletEnv env.ApateletEnvironment, kubernetesPort, metricsPort int, readyCh chan<- struct{}) error {
	log.Println("Starting Apatelet")

	// Retrieving connection information
	connectionInfo := service.NewConnectionInfo(apateletEnv.ControlPlaneAddress, apateletEnv.ControlPlanePort, false)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create stop channels
	stop := make(chan os.Signal, 1)
	stopInformer := make(chan struct{})

	// Join the apate cluster
	config, res, startTime, err := joinApateCluster(ctx, connectionInfo, apateletEnv.ListenPort)
	if err != nil {
		return errors.Wrap(err, "failed to join apate cluster")
	}

	// Write kubeconfig if it doesn't exist
	writeKubeConfig(apateletEnv, config)

	// Create store
	st := store.NewStore()

	// Create scheduler
	sch := createScheduler(ctx, st)

	// Create crd informers
	err = createInformers(config, st, stopInformer, sch, res)
	if err != nil {
		return errors.Wrap(err, "failed to create informers")
	}

	// Setup health status
	hc, err := startHealth(ctx, connectionInfo, res.UUID, stop)
	if err != nil {
		return errors.Wrap(err, "failed to start health client")
	}

	// Start the Apatelet
	nc, err := vkProvider.CreateProvider(ctx, &apateletEnv, res, kubernetesPort, metricsPort, &st)
	if err != nil {
		return errors.Wrap(err, "failed to create provider")
	}

	// Create virtual kubelet
	log.Println("Joining kubernetes cluster")
	errch := make(chan error)
	go func() {
		if err = nc.Run(ctx); err != nil {
			hc.SetStatus(healthpb.Status_UNHEALTHY)
			errch <- errors.Wrap(err, "failed to run node controller")
		}
	}()

	// Start gRPC server
	server, err := createGRPC(apateletEnv.ListenPort, &st, &sch, apateletEnv.ListenAddress, stop)
	if err != nil {
		return errors.Wrap(err, "failed to set up GRPC endpoints")
	}

	// Update status
	hc.SetStatus(healthpb.Status_HEALTHY)
	log.Printf("now accepting requests on %s:%d\n", server.Conn.Address, server.Conn.Port)
	log.Printf("now listening on :%d for kube api and :%d for metrics", kubernetesPort, metricsPort)

	// Handle stop
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Start serving request
	go func() {
		if err := server.Serve(); err != nil {
			errch <- errors.Wrap(err, "apatelet server failed")
		}
	}()

	// Start the scheduler if a scenario is already running
	if startTime >= 0 {
		sch.StartScheduler(startTime)
	}

	readyCh <- struct{}{}

	// Stop the server on signal or error
	select {
	case err := <-errch:
		err = errors.Wrap(err, "apatelet stopped because of an error")
		log.Println(err)
		return err
	case <-stop:
		stopInformer <- struct{}{}
		if err := shutdown(ctx, server, connectionInfo, res.UUID.String()); err != nil {
			log.Println(err)
		}
		return nil
	}
}

func createScheduler(ctx context.Context, st store.Store) scheduler.Scheduler {
	sch := scheduler.New(ctx, &st)
	go func() {
		ech := sch.EnableScheduler()

		for {
			select {
			case <-ctx.Done():
				return
			case err := <-ech:
				fmt.Printf("error while scheduling task occurred: %v\n", err)
			}
		}
	}()

	return sch
}

func writeKubeConfig(apateletEnv env.ApateletEnvironment, config *kubeconfig.KubeConfig) {
	once.Do(func() {
		kubeConfigLocation := apateletEnv.KubeConfigLocation
		_, err := os.Stat(kubeConfigLocation)
		if os.IsNotExist(err) {
			return
		} else if err != nil {
			panic(errors.Wrap(err, "error while reading kubeconfig file"))
		}

		err = ioutil.WriteFile(kubeConfigLocation, config.Bytes, 0600)
		if err != nil {
			panic(errors.Wrap(err, "error while writing kubeconfig to file"))
		}
	})
}

func createInformers(config *kubeconfig.KubeConfig, st store.Store, stopInformer chan struct{}, sch scheduler.Scheduler, res *scenario.NodeResources) error {
	err := crdPod.CreatePodInformer(config, &st, stopInformer, sch.WakeScheduler)
	if err != nil {
		return errors.Wrap(err, "failed creating crd pod informer")
	}

	err = crdNode.CreateNodeInformer(config, &st, res.Selector, stopInformer, sch.WakeScheduler)
	if err != nil {
		return errors.Wrap(err, "failed creating crd node informer")
	}

	return nil
}

func shutdown(ctx context.Context, server *service.GRPCServer, connectionInfo *service.ConnectionInfo, uuid string) error {
	log.Println("Stopping Apatelet")

	log.Println("Stopping API")
	server.Server.Stop()

	log.Println("Leaving clusters (apate & k8s)")

	client, err := controlplane.GetClusterOperationClient(connectionInfo)
	if err != nil {
		return errors.Wrap(err, "failed to get cluster operation client")
	}
	defer func() {
		err := client.Conn.Close()
		if err != nil {
			log.Printf("could not close connection: %v\n", err)
		}
	}()

	if err := client.LeaveCluster(ctx, uuid); err != nil {
		log.Printf("An error occurred while leaving the clusters (apate & k8s): %v\n", err)
	}

	log.Println("Stopped Apatelet")

	return nil
}

func startHealth(ctx context.Context, connectionInfo *service.ConnectionInfo, uuid uuid.UUID, stop chan<- os.Signal) (*health.Client, error) {
	hc, err := health.GetClient(connectionInfo, uuid.String())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get health client")
	}
	hc.SetStatus(healthpb.Status_UNKNOWN)
	var retries int32 = 3
	hc.StartStream(ctx, func(err error) {
		if atomic.LoadInt32(&retries) < 1 {
			// Stop after retries amount of errors
			stop <- syscall.SIGTERM
			return
		}
		log.Println(err)
		atomic.AddInt32(&retries, -1)
	})
	return hc, nil
}

func joinApateCluster(ctx context.Context, connectionInfo *service.ConnectionInfo, listenPort int) (*kubeconfig.KubeConfig, *scenario.NodeResources, int64, error) {
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

	cfg, res, startTime, err := client.JoinCluster(ctx, listenPort)

	if err != nil {
		return nil, nil, -1, errors.Wrap(err, "failed to join cluster")
	}

	log.Printf("Joined apate cluster with resources: %v", res)

	return cfg, res, startTime, nil
}

func createGRPC(listenPort int, store *store.Store, sch *scheduler.Scheduler, listenAddress string, stopChannel chan<- os.Signal) (*service.GRPCServer, error) {
	// Connection settings
	connectionInfo := service.NewConnectionInfo(listenAddress, listenPort, false)

	// Create gRPC server
	server, err := service.NewGRPCServer(connectionInfo)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new GRPC server")
	}

	// Add services
	vkService.RegisterScenarioService(server, store, sch)
	vkService.RegisterApateletService(server, stopChannel)

	return server, nil
}
