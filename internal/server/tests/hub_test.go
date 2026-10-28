package tests

import (
	"Chat/internal/server"
	"testing"
)

func TestHubRejectsTakenNickname(t *testing.T) {
	hub := server.NewHub()

	first := &server.Session{}
	second := &server.Session{}

	hub.Add(first)
	hub.Add(second)

	err := hub.Rename(first, "alex")
	if err != nil {
		t.Fatalf("unable to set first nickname: %v", err)
	}

	previousName := second.GetNickname()

	err = hub.Rename(second, "alex")
	if err == nil {
		t.Fatal("expected an error when taking an occupied nickname")
	}

	if second.GetNickname() != previousName {
		t.Fatal("second client's nickname changed after rejection")
	}

	if first.GetNickname() != "alex" {
		t.Fatal("first client lost its nickname")
	}
}
