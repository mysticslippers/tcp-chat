package commands

import "fmt"

type NicknameCommand struct {
	Name   string
	Target *string
}

func (command NicknameCommand) Execute() string {
	*command.Target = command.Name

	return fmt.Sprintf("Nickname changed to %s\n", command.Name)
}
