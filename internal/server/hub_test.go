package server

import "testing"

func TestHubRejectsTakenNickname(t *testing.T) {
	hub := NewHub()

	first := &Session{hub: hub}
	second := &Session{hub: hub}

	hub.Add(first)
	hub.Add(second)

	err := first.SetNickname("alex")
	if err != nil {
		t.Fatalf("unable to set first nickname: %v", err)
	}

	previousName := second.GetNickname()

	err = second.SetNickname("alex")
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
