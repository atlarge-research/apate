package run

import (
	"context"
	"log"
	"os"
	"sync/atomic"
	"syscall"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	healthpb "github.com/atlarge-research/opendc-emulate-kubernetes/api/health"
	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/service"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/clients/health"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/scheduler"
	vkService "github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/services"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
)

func createGRPC(store *store.Store, sch *scheduler.Scheduler, listenAddress string, stopChannel chan<- os.Signal) (*service.GRPCServer, error) {
	// Connection settings
	connectionInfo := service.NewConnectionInfo(listenAddress, 0, false)

	// Create gRPC server
	server, err := service.NewGRPCServer(true, connectionInfo)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new GRPC server")
	}

	// Add services
	vkService.RegisterScenarioService(server, store, sch)
	vkService.RegisterApateletService(server, stopChannel)

	return server, nil
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
			log.Printf("stopping apatelet because of health stream failure")
			stop <- syscall.SIGTERM
			return
		}
		log.Println(err)
		atomic.AddInt32(&retries, -1)
	})
	return hc, nil
}
