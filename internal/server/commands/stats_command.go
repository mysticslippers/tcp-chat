package commands

import "fmt"

type StatsCommand struct {
	MessageCount int
}

func (command StatsCommand) Execute() string {
	return fmt.Sprintf("Messages sent: %d\n", command.MessageCount)
}
