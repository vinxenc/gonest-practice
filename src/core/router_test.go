package core

import (
	"fmt"
	"strings"
	"testing"

	"github.com/danielgtaylor/huma/v2"
)

type fakeRoute struct{}

func (fakeRoute) RegisterRoutes(huma.API) {}

func newFakeRoute() *fakeRoute { return &fakeRoute{} }

func newRouteWithError() (*fakeRoute, error) { return &fakeRoute{}, nil }

func newNonRoute() *struct{} { return &struct{}{} }

func newRouteWithExtra() (*fakeRoute, *struct{}) { return &fakeRoute{}, &struct{}{} }

func newRouteSecondPosition() (*struct{}, *fakeRoute) { return &struct{}{}, &fakeRoute{} }

func TestAsRoute_AcceptsRouteConstructor(t *testing.T) {
	if got := AsRoute(newFakeRoute); got == nil {
		t.Fatal("AsRoute returned nil for a valid Route constructor")
	}
}

func TestAsRoute_AcceptsRouteWithErrorConstructor(t *testing.T) {
	if got := AsRoute(newRouteWithError); got == nil {
		t.Fatal("AsRoute returned nil for a (Route, error) constructor")
	}
}

func TestAsRoute_PanicsOnNonFunction(t *testing.T) {
	defer func() {
		if msg := fmt.Sprint(recover()); !strings.Contains(msg, "constructor must be a function") {
			t.Fatalf("unexpected panic: %v", msg)
		}
	}()
	AsRoute("not a function")
}

func TestAsRoute_PanicsWhenResultIsNotRoute(t *testing.T) {
	defer func() {
		if msg := fmt.Sprint(recover()); !strings.Contains(msg, "must return Route or (Route, error)") {
			t.Fatalf("unexpected panic: %v", msg)
		}
	}()
	AsRoute(newNonRoute)
}

func TestAsRoute_PanicsOnExtraNonErrorResult(t *testing.T) {
	defer func() {
		if msg := fmt.Sprint(recover()); !strings.Contains(msg, "must return Route or (Route, error)") {
			t.Fatalf("unexpected panic: %v", msg)
		}
	}()
	AsRoute(newRouteWithExtra)
}

func TestAsRoute_PanicsWhenRouteNotFirst(t *testing.T) {
	defer func() {
		if msg := fmt.Sprint(recover()); !strings.Contains(msg, "must return Route or (Route, error)") {
			t.Fatalf("unexpected panic: %v", msg)
		}
	}()
	AsRoute(newRouteSecondPosition)
}
