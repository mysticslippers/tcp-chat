package commands

import "fmt"

type WhereCommand struct {
	RoomName string
}

func (command WhereCommand) Execute() string {
	return fmt.Sprintf("Current room: %s\n", command.RoomName)
}
