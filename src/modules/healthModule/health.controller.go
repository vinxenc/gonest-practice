package healthModule

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// Controller registers and handles the health module's HTTP routes.
type Controller struct {
	service *Service
}

// HealthController constructs a health Controller with its service (fx provider).
func HealthController(service *Service) *Controller {
	return &Controller{service: service}
}

// RegisterRoutes wires the health endpoints onto the given Huma API.
func (c *Controller) RegisterRoutes(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "get-health",
		Method:      http.MethodGet,
		Path:        "/health",
		Summary:     "Health check",
		Description: "Returns the current health status of the service.",
		Tags:        []string{"Health"},
	}, func(ctx context.Context, input *struct{}) (*HealthOutput, error) {
		resp := &HealthOutput{}
		resp.Body.Status = c.service.Check()
		return resp, nil
	})
}
