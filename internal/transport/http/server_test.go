package http_test

import (
	"context"
	"errors"
	"net/http"
	"syscall"
	"testing"
	"time"

	transporthttp "github.com/cybre/salesforge-assignment/internal/transport/http"
)

func TestServer_Start_Success(t *testing.T) {
	// Create a new Server instance
	s := transporthttp.Server{}

	// Create a context with a cancel function
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the server in a separate goroutine
	go func() {
		err := s.Start(ctx, "3001")
		if err != nil {
			t.Errorf("server returned an error: %v", err)
		}
	}()

	// Wait for the server to start
	time.Sleep(100 * time.Millisecond)

	// Send a request to the server
	resp, err := http.Get("http://localhost:3001/health")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, but got %d", resp.StatusCode)
	}

	// Cancel the context to stop the server
	cancel()

	// Wait for the server to stop
	time.Sleep(100 * time.Millisecond)

	// Send a request to the server after it has stopped
	_, err = http.Get("http://localhost:3001/health")
	if err == nil {
		t.Error("expected an error, but got nil")
	}
	if !errors.Is(err, syscall.ECONNREFUSED) {
		t.Errorf("unexpected error: %v", err)
	}
}
