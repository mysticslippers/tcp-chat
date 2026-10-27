package commands

type MessageCommand interface {
	Send(message string) error
}
