package main

import tea "charm.land/bubbletea/v2"

func (model chatModel) View() tea.View {
	status := "Подключение…"
	if model.connection != nil {
		status = "Подключено"
	} else if model.networkErr != nil {
		status = "Переподключение…"
	}

	content := "TCP Chat | " + status + "\n" +
		model.viewport.View() + "\n" +
		model.input.View() + "\n" +
		"Enter: отправить | PgUp/PgDn: прокрутка | Ctrl+C: выход"

	view := tea.NewView(content)
	view.AltScreen = true

	return view
}
