package tests

import (
	"Chat/internal/server"
	"fmt"
	"testing"
)

func TestHubHistoryKeepsLast20Messages(t *testing.T) {
	hub := server.NewHub()

	for number := 1; number <= 25; number++ {
		hub.SaveMessage(fmt.Sprintf("message-%d\n", number))
	}

	history := hub.GetHistory()

	if len(history) != 20 {
		t.Fatalf("expected 20 messages, got %d", len(history))
	}

	for index, message := range history {
		expected := fmt.Sprintf("message-%d\n", index+6)

		if message != expected {
			t.Fatalf(
				"history[%d]: expected %q, got %q",
				index, expected, message,
			)
		}
	}
}

func TestHubHistoryReturnsCopy(t *testing.T) {
	hub := server.NewHub()
	hub.SaveMessage("original\n")

	history := hub.GetHistory()

	if len(history) != 1 {
		t.Fatalf("expected one message, got %d", len(history))
	}

	history[0] = "modified\n"

	actual := hub.GetHistory()

	if len(actual) != 1 {
		t.Fatalf("expected one stored message, got %d", len(actual))
	}

	if actual[0] != "original\n" {
		t.Fatalf("internal history was modified: %q", actual[0])
	}
}
