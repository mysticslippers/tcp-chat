package main

import (
	"errors"
	"log"
	"net"

	tea "charm.land/bubbletea/v2"
)

func main() {
	program := tea.NewProgram(newChatModel(nil))

	_, err := program.Run()
	if err != nil {
		log.Fatalf("Client error: %s", err)
	}
}

func run(address string) error {
	connection, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}

	defer func(connection net.Conn) {
		err := connection.Close()
		if err != nil && !errors.Is(err, net.ErrClosed) {
			log.Printf("Error closing connection: %s", err)
		}
	}(connection)

	program := tea.NewProgram(newChatModel(connection))

	result, err := program.Run()
	if err != nil {
		return err
	}

	return result.(chatModel).networkErr
}
