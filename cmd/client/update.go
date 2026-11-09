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

			model.messages = append(model.messages, "Вы: "+text)
			model.viewport.SetContent(strings.Join(model.messages, "\n"))
			model.viewport.GotoBottom()
			model.input.Reset()

			return model, nil

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
	}

	var command tea.Cmd
	model.input, command = model.input.Update(event)

	return model, command
}
