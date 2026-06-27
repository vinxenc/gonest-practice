package core

import (
	"context"
	"net"
	"strconv"
	"testing"
	"time"

	"gonest-practice/src/config"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

// fakeLifecycle records appended hooks so a test can invoke them directly,
// standing in for fx's real lifecycle.
type fakeLifecycle struct{ hooks []fx.Hook }

func (l *fakeLifecycle) Append(h fx.Hook) { l.hooks = append(l.hooks, h) }

// fakeShutdowner records whether Shutdown was asked for, standing in for fx's
// Shutdowner.
type fakeShutdowner struct {
	called bool
	opts   []fx.ShutdownOption
}

func (s *fakeShutdowner) Shutdown(opts ...fx.ShutdownOption) error {
	s.called = true
	s.opts = opts
	return nil
}

func TestNewFiber(t *testing.T) {
	if app := NewFiber(); app == nil {
		t.Fatal("NewFiber returned nil")
	}
}

func TestNewHumaAPI(t *testing.T) {
	if api := NewHumaAPI(NewFiber()); api == nil {
		t.Fatal("NewHumaAPI returned nil")
	}
}

// freePort asks the OS for an unused TCP port and returns it, closing the
// temporary listener so a caller can bind it.
func freePort(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("reserving a free port: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	if err := ln.Close(); err != nil {
		t.Fatalf("closing reservation listener: %v", err)
	}
	return port
}

// TestStartServer_StartStop drives the lifecycle hook startServer registers: the
// OnStart hook must bind successfully and OnStop must shut the server down
// cleanly.
func TestStartServer_StartStop(t *testing.T) {
	port := freePort(t)
	lc := &fakeLifecycle{}
	startServer(lc, fiber.New(), &fakeShutdowner{}, &config.Settings{Port: port})

	if len(lc.hooks) != 1 {
		t.Fatalf("startServer appended %d hooks, want 1", len(lc.hooks))
	}
	h := lc.hooks[0]
	ctx := context.Background()

	if err := h.OnStart(ctx); err != nil {
		t.Fatalf("OnStart returned error: %v", err)
	}

	// OnStart serves on the listener in a goroutine, so wait until the server is
	// actually accepting connections before stopping it. Without this, OnStop can
	// race the serving goroutine and leave it running after shutdown.
	addr := "127.0.0.1:" + strconv.Itoa(port)
	if !waitListening(addr, 5*time.Second) {
		t.Fatalf("server never started listening on %s", addr)
	}

	if err := h.OnStop(ctx); err != nil {
		t.Fatalf("OnStop returned error: %v", err)
	}
}

// waitListening reports whether addr accepts a TCP connection before the timeout.
func waitListening(addr string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", addr, 100*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			return true
		}
		time.Sleep(20 * time.Millisecond)
	}
	return false
}

// TestStartServer_BindError verifies a failed bind (port already in use)
// propagates out of OnStart so fx can fail startup, rather than crashing a
// goroutine.
func TestStartServer_BindError(t *testing.T) {
	// Occupy a port on all interfaces so startServer's ":port" bind conflicts.
	busy, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("occupying a port: %v", err)
	}
	defer func() { _ = busy.Close() }()
	port := busy.Addr().(*net.TCPAddr).Port

	lc := &fakeLifecycle{}
	startServer(lc, fiber.New(), &fakeShutdowner{}, &config.Settings{Port: port})

	if len(lc.hooks) != 1 {
		t.Fatalf("startServer appended %d hooks, want 1", len(lc.hooks))
	}
	if err := lc.hooks[0].OnStart(context.Background()); err == nil {
		t.Fatal("OnStart on an occupied port = nil error, want bind error")
	}
}

// TestInitServer verifies initServer registers every controller's routes and
// wires the server's start/stop hook onto the lifecycle.
func TestInitServer(t *testing.T) {
	rc := &recordingController{}
	lc := &fakeLifecycle{}
	app := fiber.New()

	initServer(serverParams{
		Lifecycle:   lc,
		Shutdowner:  &fakeShutdowner{},
		App:         app,
		API:         NewHumaAPI(app),
		Settings:    &config.Settings{Port: freePort(t)},
		Controllers: []Controller{rc},
	})

	if rc.calls != 1 {
		t.Fatalf("initServer registered controller %d times, want 1", rc.calls)
	}
	if len(lc.hooks) != 1 {
		t.Fatalf("initServer appended %d lifecycle hooks, want 1", len(lc.hooks))
	}
}

// TestServer_StartStop builds the full application via Server and runs it through
// a real fx start/stop, exercising the composition root end to end: config load,
// Fiber + Huma providers, controller-group collection, route registration and
// the server lifecycle.
func TestServer_StartStop(t *testing.T) {
	t.Setenv("PORT", strconv.Itoa(freePort(t)))

	rc := &recordingController{}
	app := Server(
		fx.Provide(AsController(func() *recordingController { return rc })),
	)

	ctx := context.Background()
	if err := app.Start(ctx); err != nil {
		t.Fatalf("app.Start returned error: %v", err)
	}
	if rc.calls != 1 {
		t.Fatalf("controller registered %d times during startup, want 1", rc.calls)
	}
	if err := app.Stop(ctx); err != nil {
		t.Fatalf("app.Stop returned error: %v", err)
	}
}
