package commands

import "fmt"

type NicknameSetter interface {
	SetNickname(name string) error
}

type NicknameCommand struct {
	Name   string
	Target NicknameSetter
}

func (command NicknameCommand) Execute() string {
	err := command.Target.SetNickname(command.Name)
	if err != nil {
		return fmt.Sprintf("Unable to change nickname: %s\n", err)
	}

	return fmt.Sprintf("Nickname changed to %s\n", command.Name)
}
