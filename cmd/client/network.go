package main

import (
	"encoding/json"

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
