package healthModule

import (
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"gonest-practice/src/core"
)

// TestHealthModule_ProvidesControllerToGroup verifies the module wiring: building
// the HealthModule must contribute a controller into the "controllers" value
// group whose routes register the /health endpoint.
func TestHealthModule_ProvidesControllerToGroup(t *testing.T) {
	var controllers []core.Controller
	app := fxtest.New(t,
		HealthModule,
		fx.Invoke(fx.Annotate(
			func(cs []core.Controller) { controllers = cs },
			fx.ParamTags(`group:"controllers"`),
		)),
	)
	defer app.RequireStop()
	app.RequireStart()

	if len(controllers) != 1 {
		t.Fatalf("HealthModule contributed %d controllers, want 1", len(controllers))
	}

	// The contributed controller must register the /health route.
	_, api := humatest.New(t)
	controllers[0].RegisterRoutes(api)
	if resp := api.Get("/health"); resp.Code != http.StatusOK {
		t.Fatalf("module controller GET /health = %d, want %d", resp.Code, http.StatusOK)
	}
}
