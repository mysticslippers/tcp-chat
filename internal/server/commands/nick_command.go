package commands

import "fmt"

type NicknameSetter interface {
	SetNickname(name string)
}

type NicknameCommand struct {
	Name   string
	Target NicknameSetter
}

func (command NicknameCommand) Execute() string {
	command.Target.SetNickname(command.Name)

	return fmt.Sprintf("Nickname changed to %s\n", command.Name)
}
