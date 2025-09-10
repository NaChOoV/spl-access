package server

import (
	"context"
	"spl-access/src/config"
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
	authMiddleware *middleware.AuthMiddleware,
	accessWebsocket websocket.AccessWb,
	config *config.EnvironmentConfig,
) {
	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	// Setup routes
	app.Get("/health", mainController.Health)
	app.Get("/api/access", accessController.GetObfuscateAccess)
	app.Get("/api/access/complete", authMiddleware.ValidateAuthHeader, accessController.GetAccessComplete)
	app.Post("/api/access", authMiddleware.ValidateAuthHeader, accessController.UpdateOrCreateAccess)
	app.Get("/ws/access", func(c *fiber.Ctx) error { return accessWebsocket.Upgrade(c) })

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			port := config.Port
			go app.Listen(":" + port)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return app.Shutdown()
		},
	})
}
