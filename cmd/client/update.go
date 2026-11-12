package main

import (
	"strings"

	tea "charm.land/bubbletea/v2"
)

func (model chatModel) Update(event tea.Msg) (tea.Model, tea.Cmd) {
	switch event := event.(type) {
	case tea.WindowSizeMsg:
		model.viewport.SetWidth(max(1, event.Width))
		model.viewport.SetHeight(max(1, event.Height-3))
		model.input.SetWidth(max(1, event.Width-2))
		model.viewport.GotoBottom()

		return model, nil

	case tea.KeyPressMsg:
		switch event.String() {
		case "ctrl+c":
			return model, tea.Quit

		case "enter":
			text := model.input.Value()

			if strings.TrimSpace(text) == "/quit" {
				return model, tea.Quit
			}

			if model.connection == nil || strings.TrimSpace(text) == "" {
				return model, nil
			}

			model.input.Reset()
			model.input.SetSuggestions(commandSuggestions)

			if model.sending {
				model.pending = append(model.pending, text)
				return model, nil
			}

			model.sending = true
			return model, sendMessage(model.connection, text)

		case "pgup":
			model.viewport.PageUp()
			return model, nil

		case "pgdown":
			model.viewport.PageDown()
			if model.viewport.AtBottom() {
				model.unreadCount = 0
			}
			return model, nil

		case "end":
			model.viewport.GotoBottom()
			model.unreadCount = 0
			return model, nil
		}

	case serverMessageMsg:
		if event.connection != model.connection {
			return model, nil
		}

		atBottom := model.viewport.AtBottom()

		text := strings.TrimRight(event.text, "\r\n")
		model.messages = append(model.messages, text)
		model.viewport.SetContent(strings.Join(model.messages, "\n"))

		if atBottom {
			model.viewport.GotoBottom()
		} else {
			model.unreadCount++
		}

		return model, readMessage(model.connection)

	case messageSentMsg:
		if event.connection != model.connection {
			return model, nil
		}

		if len(model.pending) == 0 {
			model.sending = false
			return model, nil
		}

		text := model.pending[0]
		model.pending = model.pending[1:]

		return model, sendMessage(model.connection, text)

	case connectResultMsg:
		if event.err != nil {
			model.networkErr = event.err
			return model, connect(model.ctx, model.address, reconnectDelay)
		}

		model.connection = event.connection
		model.networkErr = nil
		return model, readMessage(model.connection)

	case networkErrorMsg:
		if event.connection != model.connection {
			return model, nil
		}

		model.connection.cancel()
		model.connection = nil
		model.networkErr = event.err
		model.pending = nil
		model.sending = false

		return model, connect(model.ctx, model.address, 0)
	}

	previousText := model.input.Value()

	var command tea.Cmd
	model.input, command = model.input.Update(event)

	if model.input.Value() != previousText {
		model.input.SetSuggestions(commandSuggestions)
	}

	return model, command
}
