package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"time"

	tea "charm.land/bubbletea/v2"
)

const reconnectDelay = 2 * time.Second

type serverConnection struct {
	socket  net.Conn
	decoder *json.Decoder
	cancel  context.CancelFunc
}

type connectResultMsg struct {
	connection *serverConnection
	err        error
}

type serverMessageMsg string
type messageSentMsg struct{}
type networkErrorMsg struct {
	err error
}

func connect(ctx context.Context, address string, delay time.Duration) tea.Cmd {
	return func() tea.Msg {
		timer := time.NewTimer(delay)
		defer timer.Stop()

		select {
		case <-ctx.Done():
			return nil
		case <-timer.C:
		}

		dialer := net.Dialer{Timeout: 3 * time.Second}
		socket, err := dialer.DialContext(ctx, "tcp", address)
		if err != nil {
			return connectResultMsg{err: err}
		}

		connectionCtx, cancel := context.WithCancel(ctx)
		context.AfterFunc(connectionCtx, func() {
			_ = socket.Close()
		})

		return connectResultMsg{
			connection: &serverConnection{
				socket:  socket,
				decoder: json.NewDecoder(socket),
				cancel:  cancel,
			},
		}
	}
}

func readMessage(decoder *json.Decoder) tea.Cmd {
	return func() tea.Msg {
		var reply string

		err := decoder.Decode(&reply)
		if err != nil {
			return networkErrorMsg{err}
		}

		return serverMessageMsg(reply)
	}
}

func sendMessage(connection net.Conn, text string) tea.Cmd {
	return func() tea.Msg {
		err := connection.SetWriteDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			return networkErrorMsg{err}
		}

		_, err = fmt.Fprintln(connection, text)
		if err != nil {
			return networkErrorMsg{err}
		}

		return messageSentMsg{}
	}
}
