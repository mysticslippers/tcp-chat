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

func TestHubReusesFreedGuestNickname(t *testing.T) {
	hub := server.NewHub()

	first := &server.Session{}
	second := &server.Session{}

	hub.Add(first)
	hub.Add(second)

	if first.GetNickname() != "guest-1" {
		t.Fatalf("expected guest-1, got %q", first.GetNickname())
	}

	if second.GetNickname() != "guest-2" {
		t.Fatalf("expected guest-2, got %q", second.GetNickname())
	}

	hub.Remove(first)

	third := &server.Session{}
	hub.Add(third)

	if third.GetNickname() != "guest-1" {
		t.Fatalf(
			"expected freed nickname guest-1, got %q",
			third.GetNickname(),
		)
	}

	if second.GetNickname() != "guest-2" {
		t.Fatal("existing client's nickname changed")
	}
}

func TestHubConcurrentRename(t *testing.T) {
	hub := server.NewHub()

	first := &server.Session{}
	second := &server.Session{}

	hub.Add(first)
	hub.Add(second)

	start := make(chan struct{})
	results := make(chan error, 2)

	rename := func(session *server.Session) {
		<-start
		results <- hub.Rename(session, "alex")
	}

	go rename(first)
	go rename(second)

	close(start)

	successes := 0

	for attempt := 0; attempt < 2; attempt++ {
		err := <-results
		if err == nil {
			successes++
		}
	}

	if successes != 1 {
		t.Fatalf("expected one successful rename, got %d", successes)
	}

	if first.GetNickname() == second.GetNickname() {
		t.Fatal("clients ended up with identical nicknames")
	}

	if hub.FindByNickname("alex") == nil {
		t.Fatal("nickname alex has no owner")
	}
}
