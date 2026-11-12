package main

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
)

func (model chatModel) View() tea.View {
	status := "Подключение…"
	if model.connection != nil {
		status = "Подключено"
	} else if model.networkErr != nil {
		status = "Переподключение…"
	}

	footer := "Enter: отправить | PgUp/PgDn: прокрутка | Ctrl+C: выход"
	if model.unreadCount > 0 {
		footer = fmt.Sprintf(
			"Новых сообщений: %d | End: к последнему",
			model.unreadCount,
		)
	}

	content := "TCP Chat | " + status + "\n" +
		model.viewport.View() + "\n" +
		model.input.View() + "\n" +
		footer

	view := tea.NewView(content)
	view.AltScreen = true

	return view
}
