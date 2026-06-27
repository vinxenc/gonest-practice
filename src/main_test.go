package main

import (
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"testing"
	"time"
)

// TestMain_StartsAndShutsDownOnSignal runs main() — the real composition root and
// server — in a goroutine, waits until it is accepting connections, then sends
// the process SIGTERM so fx shuts it down gracefully and main() returns.
//
// A SIGTERM that arrives before fx has installed its own signal handler would
// terminate the test binary by default disposition, so we register our own
// handler first (keeping the process alive) and resend the signal until main()
// returns, guaranteeing fx receives one after it starts listening.
func TestMain_StartsAndShutsDownOnSignal(t *testing.T) {
	// Reserve a free port for the server so startup does not depend on 3000
	// being available.
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("reserving a free port: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	if err := ln.Close(); err != nil {
		t.Fatalf("closing reservation listener: %v", err)
	}
	t.Setenv("PORT", strconv.Itoa(port))

	// Keep this test process alive against the SIGTERM we will send to fx.
	keepAlive := make(chan os.Signal, 1)
	signal.Notify(keepAlive, syscall.SIGTERM)
	defer signal.Stop(keepAlive)
	go func() {
		for range keepAlive {
		}
	}()

	done := make(chan struct{})
	go func() {
		defer close(done)
		main()
	}()

	// Wait until the server is accepting connections, which means main()'s fx
	// app has started and is about to wait for a signal.
	addr := "127.0.0.1:" + strconv.Itoa(port)
	if !waitListening(addr, 5*time.Second) {
		t.Fatalf("server never started listening on %s", addr)
	}

	// Resend SIGTERM until main() returns, so fx catches one once it is waiting.
	deadline := time.After(10 * time.Second)
	for {
		if err := syscall.Kill(syscall.Getpid(), syscall.SIGTERM); err != nil {
			t.Fatalf("sending SIGTERM: %v", err)
		}
		select {
		case <-done:
			return
		case <-time.After(200 * time.Millisecond):
		case <-deadline:
			t.Fatal("main() did not return after SIGTERM")
		}
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
