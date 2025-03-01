package server

import (
	"context"
	"spl-access/src/controller"
	"spl-access/src/middleware"
	"spl-access/src/websocket"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.uber.org/fx"
)

func CreateFiberServer(
	lc fx.Lifecycle,
	accessController *controller.AccessController,
	mainController *controller.MainController,
	websocketController *websocket.WebsocketController,
	authMiddleware *middleware.AuthMiddleware,
) {
	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	// Setup routes
	app.Get("/health", mainController.Health)
	app.Get("/api/access", accessController.GetObfuscateAccess)
	app.Get("/api/access/complete", authMiddleware.ValidateAuthHeader, accessController.GetAccess)
	app.Post("/api/access", authMiddleware.ValidateAuthHeader, accessController.UpdateOrCreateAccess)
	app.Get("/ws/access", websocketController.WsAccess)

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			// TODO: switch the port to an env variable
			go app.Listen(":8000")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return app.Shutdown()
		},
	})
}
