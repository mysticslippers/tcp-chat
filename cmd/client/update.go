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

			if strings.TrimSpace(text) == "" {
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
			return model, nil
		}

	case serverMessageMsg:
		atBottom := model.viewport.AtBottom()

		text := strings.TrimRight(string(event), "\r\n")
		model.messages = append(model.messages, text)
		model.viewport.SetContent(strings.Join(model.messages, "\n"))

		if atBottom {
			model.viewport.GotoBottom()
		}

		return model, readMessage(model.decoder)

	case messageSentMsg:
		if len(model.pending) == 0 {
			model.sending = false
			return model, nil
		}

		text := model.pending[0]
		model.pending = model.pending[1:]

		return model, sendMessage(model.connection, text)

	case networkErrorMsg:
		model.networkErr = event.err

		return model, tea.Quit
	}

	previousText := model.input.Value()

	var command tea.Cmd
	model.input, command = model.input.Update(event)

	if model.input.Value() != previousText {
		model.input.SetSuggestions(commandSuggestions)
	}

	return model, command
}
