package core

import (
	"context"
	"fmt"
	"log"
	"net"

	"gonest-practice/src/config"

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
	cfg := huma.DefaultConfig("Gonest Practice API", "1.0.0")
	return humafiber.NewV2(app, cfg)
}

// startServer ties the Fiber app to the fx application lifecycle, listening on
// the port from the validated Settings.
func startServer(lc fx.Lifecycle, app *fiber.App, shutdowner fx.Shutdowner, settings *config.Settings) {
	addr := fmt.Sprintf(":%d", settings.Port)
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Bind synchronously so listen failures (e.g. port in use)
			// propagate through fx startup instead of crashing a goroutine.
			ln, err := net.Listen("tcp", addr)
			if err != nil {
				return err
			}
			go func() {
				// app.Listener returns nil on graceful shutdown (OnStop). Any
				// other error means the server died unexpectedly, so ask fx to
				// shut the app down rather than leaving it running without a
				// listening HTTP server.
				if err := app.Listener(ln); err != nil {
					log.Printf("fiber server stopped: %v", err)
					_ = shutdowner.Shutdown(fx.ExitCode(1))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return app.ShutdownWithContext(ctx)
		},
	})
}
