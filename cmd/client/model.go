package main

import (
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
)

type chatModel struct {
	input    textinput.Model
	viewport viewport.Model
	messages []string
}

func buildChatModel() chatModel {
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
	}
}
