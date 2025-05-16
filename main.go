package main

import (
	"spl-access/src/config"
	"spl-access/src/controller"
	"spl-access/src/db"
	"spl-access/src/middleware"
	"spl-access/src/repository"
	"spl-access/src/server"
	"spl-access/src/service"
	"spl-access/src/websocket"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		// Setup database
		fx.Provide(db.CreatePostgresConnection),
		// Setup repository
		fx.Provide(
			fx.Annotate(
				repository.NewPostgresAccess,
				fx.As(new(repository.AccessRepository)),
			),
		),
		fx.Provide(repository.NewUserRepository),
		// Setup background context
		fx.Provide(config.NewContextBackground),
		// Setup service
		fx.Provide(service.NewAccessService),
		// Setup controller
		fx.Provide(controller.NewAccessController),
		fx.Provide(controller.NewMainController),
		fx.Provide(websocket.NewWebsocketController),
		// Setup environment config
		fx.Provide(config.NewEnviromentConfig),
		fx.Provide(middleware.NewAuthMiddleware),
		fx.Invoke(server.CreateFiberServer),
		fx.Invoke(func(accessService *service.AccessService) {
			accessService.UpdateAccess()
		}),
	).Run()
}
