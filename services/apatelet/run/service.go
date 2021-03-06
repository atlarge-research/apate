package run

import (
	"context"
	"log"
	"os"
	"sync/atomic"
	"syscall"

	"github.com/atlarge-research/apate/pkg/channel"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	healthpb "github.com/atlarge-research/apate/api/health"
	"github.com/atlarge-research/apate/internal/service"
	"github.com/atlarge-research/apate/pkg/clients/health"
	"github.com/atlarge-research/apate/services/apatelet/scheduler"
	vkService "github.com/atlarge-research/apate/services/apatelet/services"
	"github.com/atlarge-research/apate/services/apatelet/store"
)

func createGRPC(store *store.Store, sch *scheduler.Scheduler, listenAddress string, listenPort int, stopCh chan<- struct{}, stopInformerCh *channel.StopChannel) (*service.GRPCServer, error) {
	// Connection settings
	connectionInfo := service.NewConnectionInfo(listenAddress, listenPort)

	// Create gRPC server
	server, err := service.NewGRPCServer(connectionInfo)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new GRPC server")
	}

	// Add services
	vkService.RegisterScenarioService(server, store, sch, stopInformerCh)
	vkService.RegisterApateletService(server, stopCh)

	return server, nil
}

func startHealth(ctx context.Context, connectionInfo *service.ConnectionInfo, uuid uuid.UUID, stop chan<- os.Signal) (*health.Client, error) {
	hc, err := health.GetClient(connectionInfo, uuid.String())
	if err != nil {
		return nil, errors.Wrap(err, "failed to get health client")
	}
	hc.SetStatus(healthpb.Status_UNKNOWN)
	var retries int32 = 3
	hc.StartStream(ctx, func(err error) bool {
		if atomic.LoadInt32(&retries) < 1 {
			// Stop after retries amount of errors
			select {
			case stop <- syscall.SIGTERM:
				log.Printf("stopping apatelet because of health stream failure")
			default:
				//
			}

			return true
		}
		log.Println(err)
		atomic.AddInt32(&retries, -1)
		return false
	})
	return hc, nil
}
