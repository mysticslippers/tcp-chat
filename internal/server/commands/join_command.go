package commands

type JoinCommand interface {
	JoinRoom(name string) error
}
