package main

import tea "charm.land/bubbletea/v2"

func (model chatModel) View() tea.View {
	content := "TCP Chat\n" +
		model.viewport.View() + "\n" +
		model.input.View() + "\n" +
		"Enter: добавить | PgUp/PgDn: прокрутка | Ctrl+C: выход"

	view := tea.NewView(content)
	view.AltScreen = true

	return view
}
