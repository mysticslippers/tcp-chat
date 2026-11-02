package commands

import "fmt"

type RoomJoiner interface {
	JoinRoom(name string) error
}

type JoinCommand struct {
	Name   string
	Target RoomJoiner
}

func (command JoinCommand) Execute() string {
	err := command.Target.JoinRoom(command.Name)
	if err != nil {
		return fmt.Sprintf("Unable to join room: %s\n", err)
	}

	return fmt.Sprintf("Joined room %s\n", command.Name)
}
