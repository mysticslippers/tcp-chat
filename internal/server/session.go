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
	"sync"
)

type Session struct {
	MessageCount int

	nickname  string
	nameMutex sync.Mutex

	encoder *json.Encoder
	mutex   sync.Mutex
}

func (session *Session) SetNickname(name string) {
	session.nameMutex.Lock()
	defer session.nameMutex.Unlock()

	session.nickname = name
}

func (session *Session) GetNickname() string {
	session.nameMutex.Lock()
	defer session.nameMutex.Unlock()

	return session.nickname
}

func (session *Session) Send(message string) error {
	session.mutex.Lock()
	defer session.mutex.Unlock()

	return session.encoder.Encode(message)
}

func Handle(connection net.Conn, stop <-chan struct{}, hub *Hub) error {
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

	session := &Session{
		nickname: "guest",
		encoder:  json.NewEncoder(connection),
	}

	err := session.Send("You have successfully connected to the server\n")
	if err != nil {
		return err
	}

	hub.Add(session)
	defer hub.Remove(session)

	reader := bufio.NewReader(connection)
	for {
		message, err := reader.ReadString('\n')
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		parts := strings.Fields(message)

		commandName := ""
		if len(parts) > 0 {
			commandName = parts[0]
		}

		switch commandName {
		case "/stats":
			var command commands.Command = commands.StatsCommand{
				MessageCount: session.MessageCount,
			}

			message = command.Execute()

		case "/help":
			var command commands.Command = commands.HelpCommand{}

			message = command.Execute()

		case "/nick":
			if len(parts) == 2 {
				var command commands.Command = commands.NicknameCommand{
					Name:   parts[1],
					Target: &session.nickname,
				}

				message = command.Execute()
			} else {
				message = "Usage: /nick <name>\n"
			}

		case "/online":
			var command commands.Command = commands.OnlineCommand{
				Count: hub.Count(),
			}

			message = command.Execute()

		default:
			session.MessageCount++
			message = session.nickname + ": " + message

			hub.Broadcast(message)
			continue
		}

		err = session.Send(message)
		if err != nil {
			return err
		}
	}
}
