package healthModule

import (
	"net/http"

	"github.com/0xfurai/gonest"
)

// Controller registers and handles the health module's HTTP routes.
type Controller struct {
	service *Service
}

// HealthController constructs a health Controller with its service (gonest
// provider).
func HealthController(service *Service) *Controller {
	return &Controller{service: service}
}

// Register wires the health endpoints onto the given router and declares their
// OpenAPI metadata (summary, tag, response schema) for the Swagger document.
func (c *Controller) Register(r gonest.Router) {
	r.Get("/health", c.check).
		Summary("Health check").
		Tags("Health").
		Response(http.StatusOK, HealthResponse{})
}

// check handles GET /health.
func (c *Controller) check(ctx gonest.Context) error {
	return ctx.JSON(http.StatusOK, HealthResponse{Status: c.service.Check()})
}
