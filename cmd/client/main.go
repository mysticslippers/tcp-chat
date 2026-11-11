package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"

	tea "charm.land/bubbletea/v2"
)

func main() {
	const address = "localhost:9000"

	err := run(address)
	if errors.Is(err, io.EOF) {
		fmt.Println("Server closed the connection")
		return
	}

	if err != nil {
		log.Fatalf("Client error: %s", err)
	}
}

func run(address string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	program := tea.NewProgram(newChatModel(ctx, address))
	_, err := program.Run()
	return err
}
