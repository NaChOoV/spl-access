package main

import (
	"fmt"
	"spl-access/src/config"
	"spl-access/src/controller"
	"spl-access/src/db"
	"spl-access/src/middleware"
	"spl-access/src/repository"
	"spl-access/src/server"
	"spl-access/src/service"
	"spl-access/src/websocket"
	"time"

	"github.com/go-co-op/gocron/v2"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		// Setup database
		fx.Provide(db.CreateBigQueryConnection),
		// Setup repository
		fx.Provide(
			fx.Annotate(
				repository.NewBigQueryAccess,
				fx.As(new(repository.AccessRepository)),
			),
		),
		fx.Provide(
			fx.Annotate(
				repository.NewBigQueryUser,
				fx.As(new(repository.UserRepository)),
			),
		),
		// Setup background context
		fx.Provide(config.NewContextBackground),
		// Setup service
		fx.Provide(service.NewAccessService),
		fx.Provide(service.NewSourceService),
		// Setup controller
		fx.Provide(controller.NewAccessController),
		fx.Provide(controller.NewMainController),
		// Setup websocket
		fx.Provide(
			fx.Annotate(
				websocket.NewAccessWebsocket,
				fx.As(new(websocket.AccessWb)),
			),
		),
		// Setup environment config
		fx.Provide(config.NewEnviromentConfig),
		fx.Provide(middleware.NewAuthMiddleware),
		fx.Invoke(server.CreateFiberServer),
		fx.Invoke(func(accessService *service.AccessService) {
			s, err := gocron.NewScheduler()
			if err != nil {
				fmt.Println("Error creating scheduler:", err)
				return
			}

			s.NewJob(
				gocron.DurationJob(5*time.Second),
				gocron.NewTask(func() {
					err := accessService.SyncAccess()
					if err != nil {
						fmt.Println("[CRON] Error syncing access:", err)
						return
					}
				}),
				gocron.WithSingletonMode(gocron.LimitModeReschedule),
			)

			s.Start()
		}),
	).Run()
}
