package core

import (
	"context"
	"log"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

// NewFiber provides the Fiber application instance.
func NewFiber() *fiber.App {
	return fiber.New()
}

// NewHumaAPI provides the Huma API bound to the Fiber app for OpenAPI generation.
func NewHumaAPI(app *fiber.App) huma.API {
	config := huma.DefaultConfig("Gonest Practice API", "1.0.0")
	return humafiber.NewV2(app, config)
}

// startServer ties the Fiber app to the fx application lifecycle.
func startServer(lc fx.Lifecycle, app *fiber.App) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := app.Listen(":3000"); err != nil {
					log.Fatal(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return app.ShutdownWithContext(ctx)
		},
	})
}
