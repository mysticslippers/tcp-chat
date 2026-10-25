package commands

import "fmt"

type OnlineCommand struct {
	Count int
}

func (command OnlineCommand) Execute() string {
	return fmt.Sprintf("Clients online: %d\n", command.Count)
}
