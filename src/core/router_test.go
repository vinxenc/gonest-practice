package core

import (
	"testing"

	"github.com/danielgtaylor/huma/v2"
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

func TestAsController_ReturnsAnnotation(t *testing.T) {
	if got := AsController(func() *recordingController { return &recordingController{} }); got == nil {
		t.Fatal("AsController returned nil for a valid controller constructor")
	}
}
