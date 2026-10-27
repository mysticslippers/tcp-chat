package commands

type MessageTarget interface {
	Send(message string) error
}

type MessageCommand struct {
	SenderName   string
	ReceiverName string
	Message      string
	Target       MessageTarget
}
