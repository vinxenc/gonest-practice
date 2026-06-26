package core

import (
	"testing"

	"github.com/danielgtaylor/huma/v2"
)

type recordingController struct{ calls int }

func (c *recordingController) RegisterRoutes(huma.API) { c.calls++ }

type fakeModule struct{ controllers []Controller }

func (m fakeModule) Controllers() []Controller { return m.controllers }

func TestRegisterRoutes_RegistersEveryControllerPerModule(t *testing.T) {
	c1 := &recordingController{}
	c2 := &recordingController{}
	c3 := &recordingController{}

	modules := []Module{
		fakeModule{controllers: []Controller{c1, c2}},
		fakeModule{controllers: []Controller{c3}},
	}

	// The fake controllers ignore the API, so a nil value is fine here.
	registerRoutes(nil, modules)

	for i, c := range []*recordingController{c1, c2, c3} {
		if c.calls != 1 {
			t.Fatalf("controller %d: got %d RegisterRoutes calls, want 1", i, c.calls)
		}
	}
}

func TestRegisterRoutes_NoModules(t *testing.T) {
	// Must not panic when there are no modules to register.
	registerRoutes(nil, nil)
}
