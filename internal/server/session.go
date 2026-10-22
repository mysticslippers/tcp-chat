package server

import (
	"Chat/internal/server/commands"
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"strings"
)

type Session struct {
	MessageCount int
	Nickname     string
}

func Handle(connection net.Conn, stop <-chan struct{}) error {
	defer func(connection net.Conn) {
		err := connection.Close()
		if err != nil && !errors.Is(err, net.ErrClosed) {
			log.Printf("Error closing connection: %s", err)
		}
	}(connection)

	done := make(chan struct{})
	defer close(done)

	go func(connection net.Conn) {
		select {
		case <-stop:
			err := connection.Close()
			if err != nil && !errors.Is(err, net.ErrClosed) {
				log.Printf("Error closing connection: %s", err)
			}

		case <-done:
		}
	}(connection)

	log.Printf("Accepted connection from %s", connection.RemoteAddr())
	encoder := json.NewEncoder(connection)

	err := encoder.Encode("You have successfully connected to the server\n")
	if err != nil {
		return err
	}

	session := &Session{
		Nickname: "guest",
	}
	reader := bufio.NewReader(connection)
	for {
		message, err := reader.ReadString('\n')
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		switch strings.TrimSpace(message) {
		case "/stats":
			var command commands.Command = commands.StatsCommand{
				MessageCount: session.MessageCount,
			}

			message = command.Execute()

		case "/help":
			var command commands.Command = commands.HelpCommand{}

			message = command.Execute()

		default:
			session.MessageCount++
		}

		err = encoder.Encode(message)
		if err != nil {
			return err
		}
	}
}
