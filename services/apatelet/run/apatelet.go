// Package run is the entry point of the actual apatelet
package run

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"

	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/scheduler"

	"github.com/pkg/errors"

	healthpb "github.com/atlarge-research/opendc-emulate-kubernetes/api/health"
	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/service"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/clients/controlplane"
	vkProvider "github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/provider"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
)

// StartApatelet starts the apatelet
func StartApatelet(apateletEnv env.ApateletEnvironment, readyCh chan<- struct{}) error {
	log.Println("Starting Apatelet")

	// Retrieving connection information
	connectionInfo := service.NewConnectionInfo(apateletEnv.ControlPlaneAddress, apateletEnv.ControlPlanePort, false)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create stop channels
	stop := make(chan os.Signal, 1)
	stopInformer := make(chan struct{})

	// Join the apate cluster
	config, res, startTime, err := joinApateCluster(ctx, connectionInfo, apateletEnv.ListenPort, apateletEnv.KubeConfigLocation)
	if err != nil {
		return errors.Wrap(err, "failed to join apate cluster")
	}

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
	nc, err := vkProvider.CreateProvider(ctx, &apateletEnv, res, apateletEnv.KubernetesPort, apateletEnv.MetricsPort, &st)
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
	server, err := createGRPC(apateletEnv.ListenPort, &st, sch, apateletEnv.ListenAddress, stop)
	if err != nil {
		return errors.Wrap(err, "failed to set up GRPC endpoints")
	}

	// Update status
	hc.SetStatus(healthpb.Status_HEALTHY)
	log.Printf("now accepting requests on %s:%d\n", server.Conn.Address, server.Conn.Port)
	log.Printf("now listening on :%d for kube api and :%d for metrics", apateletEnv.KubernetesPort, apateletEnv.MetricsPort)

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

func createScheduler(ctx context.Context, st store.Store) *scheduler.Scheduler {
	sch := scheduler.New(&st)
	go func() {
		ech := sch.EnableScheduler(ctx)

		for {
			select {
			case <-ctx.Done():
				return
			case err := <-ech:
				fmt.Printf("error while scheduling task occurred: %v\n", err)
			}
		}
	}()

	return &sch
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
