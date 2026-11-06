package main

import (
	"log"

	tea "charm.land/bubbletea/v2"
)

func main() {
	program := tea.NewProgram(newChatModel())

	_, err := program.Run()
	if err != nil {
		log.Fatalf("Client error: %s", err)
	}
}
