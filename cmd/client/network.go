package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	tea "charm.land/bubbletea/v2"
)

type serverMessageMsg string
type messageSentMsg struct{}
type networkErrorMsg struct {
	err error
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
