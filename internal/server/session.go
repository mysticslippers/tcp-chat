package server

import (
	"Chat/internal/server/commands"
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

type Session struct {
	MessageCount int
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
	_, err := fmt.Fprintln(connection, "You have successfully connected to the server")
	if err != nil {
		return err
	}

	session := &Session{}
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

		_, err = fmt.Fprint(connection, message)
		if err != nil {
			return err
		}
	}
}
