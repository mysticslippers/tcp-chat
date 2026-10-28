package tests

import (
	"Chat/internal/server"
	"errors"
	"net"
	"os"
	"testing"
	"time"
)

func TestHandleWriteTimeout(t *testing.T) {
	serverConnection, clientConnection := net.Pipe()

	defer func() {
		_ = serverConnection.Close()
		_ = clientConnection.Close()
	}()

	hub := server.NewHub()
	stop := make(chan struct{})
	result := make(chan error, 1)

	go func() {
		result <- server.Handle(serverConnection, stop, hub)
	}()

	select {
	case err := <-result:
		if !errors.Is(err, os.ErrDeadlineExceeded) {
			t.Fatalf("expected write timeout, got: %v", err)
		}

	case <-time.After(10 * time.Second):
		t.Fatal("handler did not finish within 10 seconds")
	}
}
