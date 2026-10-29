package commands

import "strings"

type History interface {
	GetHistory() []string
}

type HistoryCommand struct {
	history History
}

func NewHistoryCommand(history History) *HistoryCommand {
	return &HistoryCommand{
		history: history,
	}
}

func (command *HistoryCommand) Execute() string {
	history := command.history.GetHistory()

	if len(history) == 0 {
		return "Chat history is empty\n"
	}

	return "Chat history:\n" + strings.Join(history, "")
}
