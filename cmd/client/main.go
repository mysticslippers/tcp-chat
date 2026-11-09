package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"

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
