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
	"time"
)

const (
	writeTimeout      = 5 * time.Second
	outgoingQueueSize = 32
)

type Session struct {
	MessageCount int
	nameMutex    sync.Mutex

	hub        *Hub
	nickname   string
	done       chan struct{}
	writerDone chan struct{}
	closeOnce  sync.Once

	encoder    *json.Encoder
	connection net.Conn
	outgoing   chan string
}

func (session *Session) isClosed() bool {
	select {
	case <-session.done:
		return true
	default:
		return false
	}
}

func (session *Session) closeResources() {
	close(session.done)

	err := session.connection.Close()
	if err != nil && !errors.Is(err, net.ErrClosed) {
		log.Printf("Error closing connection: %s", err)
	}
}

func (session *Session) close() {
	session.closeOnce.Do(session.closeResources)
}

func (session *Session) writeLoop() {
	defer close(session.writerDone)

	for {
		select {
		case <-session.done:
			return

		case message := <-session.outgoing:
			err := session.write(message)
			if err != nil {
				session.close()

				if !errors.Is(err, net.ErrClosed) {
					log.Printf("Unable to send message: %s", err)
				}

				return
			}
		}

	}
}

func (session *Session) write(message string) error {
	err := session.connection.SetWriteDeadline(time.Now().Add(writeTimeout))
	if err != nil {
		return err
	}

	return session.encoder.Encode(message)
}

func (session *Session) SetNickname(name string) error {
	return session.hub.Rename(session, name)
}

func (session *Session) setNickname(name string) {
	session.nameMutex.Lock()
	defer session.nameMutex.Unlock()

	session.nickname = name
}

func (session *Session) GetNickname() string {
	session.nameMutex.Lock()
	defer session.nameMutex.Unlock()

	return session.nickname
}

func (session *Session) enqueue(message string) error {
	if session.isClosed() {
		return net.ErrClosed
	}

	select {
	case <-session.done:
		return net.ErrClosed

	case session.outgoing <- message:
		return nil

	default:
		session.close()
		return errors.New("outgoing message queue is full")
	}
}

func (session *Session) Send(message string) error {
	if session.isClosed() {
		return net.ErrClosed
	}

	select {
	case <-session.done:
		return net.ErrClosed

	case session.outgoing <- message:
		return nil

	default:
		session.close()
		return errors.New("outgoing message queue is full")
	}
}

func Handle(connection net.Conn, stop <-chan struct{}, hub *Hub) error {
	session := &Session{
		hub:        hub,
		connection: connection,
		encoder:    json.NewEncoder(connection),
		outgoing:   make(chan string, outgoingQueueSize),
		done:       make(chan struct{}),
		writerDone: make(chan struct{}),
	}

	defer session.close()

	go func(session *Session, stop <-chan struct{}) {
		select {
		case <-stop:
			session.close()
		case <-session.done:
		}
	}(session, stop)

	log.Printf("Accepted connection from %s", connection.RemoteAddr())

	var historyCommand commands.Command = commands.NewHistoryCommand(hub)

	welcome := "You have successfully connected to the server\n"
	welcome += historyCommand.Execute()

	err := session.write(welcome)
	if err != nil {
		return err
	}

	go session.writeLoop()

	defer func(session *Session) {
		session.close()
		<-session.writerDone
	}(session)

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
					Target: session,
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

		case "/who":
			var command commands.Command = commands.WhoCommand{
				Nicknames: hub.GetNicknames(),
			}

			message = command.Execute()

		case "/msg":
			if len(parts) < 3 {
				message = "Usage: /msg <nickname> <message>\n"
				break
			}

			recipient := hub.FindByNickname(parts[1])
			if recipient == nil {
				message = "User " + parts[1] + " is not online\n"
				break
			}

			text := strings.TrimSpace(message)
			text = strings.TrimSpace(strings.TrimPrefix(text, parts[0]))
			text = strings.TrimSpace(strings.TrimPrefix(text, parts[1]))

			var command commands.Command = commands.MessageCommand{
				SenderName:   session.GetNickname(),
				ReceiverName: parts[1],
				Message:      text,
				Target:       recipient,
			}

			message = command.Execute()

		case "/history":
			var command commands.Command = commands.NewHistoryCommand(hub)
			message = command.Execute()

		default:
			session.MessageCount++
			message = session.GetNickname() + ": " + message

			hub.SaveMessage(message)
			hub.Broadcast(message)
			continue
		}

		err = session.Send(message)
		if err != nil {
			return err
		}
	}
}
