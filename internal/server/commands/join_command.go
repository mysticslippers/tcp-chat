package commands

type RoomJoiner interface {
	JoinRoom(name string) error
}

type JoinCommand struct {
	Name   string
	Target RoomJoiner
}
