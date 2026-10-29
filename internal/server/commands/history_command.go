package commands

type History interface {
	GetHistory() []string
}

type HistoryCommand struct {
	history History
}
