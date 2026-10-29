package commands

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
