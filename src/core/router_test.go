package core

import (
	"testing"

	"github.com/danielgtaylor/huma/v2"
)

type fakeRoute struct{}

func (fakeRoute) RegisterRoutes(huma.API) {}

func newFakeRoute() *fakeRoute { return &fakeRoute{} }

func newNonRoute() *struct{} { return &struct{}{} }

func TestAsRoute_AcceptsRouteConstructor(t *testing.T) {
	if got := AsRoute(newFakeRoute); got == nil {
		t.Fatal("AsRoute returned nil for a valid Route constructor")
	}
}

func TestAsRoute_PanicsOnNonFunction(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("AsRoute did not panic for a non-function argument")
		}
	}()
	AsRoute("not a function")
}

func TestAsRoute_PanicsWhenResultIsNotRoute(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("AsRoute did not panic for a constructor not returning a Route")
		}
	}()
	AsRoute(newNonRoute)
}
