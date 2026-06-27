package core

import (
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

type recordingController struct{ calls int }

func (c *recordingController) RegisterRoutes(huma.API) { c.calls++ }

func TestRegisterRoutes_RegistersEveryController(t *testing.T) {
	c1 := &recordingController{}
	c2 := &recordingController{}

	// The fake controllers ignore the API, so a nil value is fine here.
	registerRoutes(nil, []Controller{c1, c2})

	for i, c := range []*recordingController{c1, c2} {
		if c.calls != 1 {
			t.Fatalf("controller %d: got %d RegisterRoutes calls, want 1", i, c.calls)
		}
	}
}

func TestRegisterRoutes_NoControllers(t *testing.T) {
	// Must not panic when there are no controllers to register.
	registerRoutes(nil, nil)
}

// TestAsController_JoinsControllersGroup verifies the Fx grouping contract: a
// constructor wrapped with AsController must be provided into the "controllers"
// value group, so a consumer of that group receives it. If AsController ever
// regressed to returning a raw constructor (no group tag), the group consumer
// would receive nothing and this test would fail.
func TestAsController_JoinsControllersGroup(t *testing.T) {
	want := &recordingController{}

	var got []Controller
	app := fxtest.New(t,
		fx.Provide(AsController(func() *recordingController { return want })),
		fx.Invoke(fx.Annotate(
			func(controllers []Controller) { got = controllers },
			fx.ParamTags(`group:"controllers"`),
		)),
	)
	defer app.RequireStop()
	app.RequireStart()

	if len(got) != 1 || got[0] != want {
		t.Fatalf(`AsController did not wire the constructor into the "controllers" group: got %d controllers, want 1`, len(got))
	}
}
