package tests

import (
	"Chat/internal/server"
	"encoding/json"
	"fmt"
	"net"
	"testing"
	"time"
)

func TestHandleQueuesMessagesInOrder(t *testing.T) {
	serverConnection, clientConnection := net.Pipe()

	hub := server.NewHub()
	stop := make(chan struct{})
	result := make(chan error, 1)

	go func() {
		result <- server.Handle(serverConnection, stop, hub)
	}()

	defer func() {
		_ = clientConnection.Close()
		_ = serverConnection.Close()

		select {
		case <-result:
		case <-time.After(3 * time.Second):
			t.Error("handler did not stop")
		}
	}()

	err := clientConnection.SetDeadline(time.Now().Add(3 * time.Second))
	if err != nil {
		t.Fatal(err)
	}

	decoder := json.NewDecoder(clientConnection)

	var welcome string
	if err := decoder.Decode(&welcome); err != nil {
		t.Fatalf("unable to read welcome: %v", err)
	}

	messages := []string{"one", "two", "three"}

	for _, message := range messages {
		_, err := fmt.Fprintln(clientConnection, message)
		if err != nil {
			t.Fatalf("unable to send message: %v", err)
		}
	}

	for _, message := range messages {
		var reply string
		if err := decoder.Decode(&reply); err != nil {
			t.Fatalf("unable to read reply: %v", err)
		}

		expected := "guest-1: " + message + "\n"
		if reply != expected {
			t.Fatalf("expected %q, got %q", expected, reply)
		}
	}
}
