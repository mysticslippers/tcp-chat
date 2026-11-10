package main

import (
	"encoding/json"
	"net"

	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
)

var commandSuggestions = []string{
	"/help",
	"/history",
	"/join",
	"/msg",
	"/nick",
	"/online",
	"/quit",
	"/rooms",
	"/stats",
	"/where",
	"/who",
}

type chatModel struct {
	input    textinput.Model
	viewport viewport.Model
	messages []string

	connection net.Conn
	decoder    *json.Decoder
	pending    []string
	sending    bool
	networkErr error
}

func newChatModel(connection net.Conn) chatModel {
	input := textinput.New()
	input.Prompt = "> "
	input.Placeholder = "Сообщение или команда... Можете ввести команду /help"
	input.SetVirtualCursor(true)
	input.Focus()

	view := viewport.New(
		viewport.WithWidth(80),
		viewport.WithHeight(20),
	)
	view.SoftWrap = true
	view.FillHeight = true

	return chatModel{
		input:    input,
		viewport: view,

		connection: connection,
		decoder:    json.NewDecoder(connection),
	}
}

func (model chatModel) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		readMessage(model.decoder),
	)
}
