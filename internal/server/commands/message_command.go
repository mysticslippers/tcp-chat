package commands

import "fmt"

type MessageTarget interface {
	Send(message string) error
}

type MessageCommand struct {
	SenderName   string
	ReceiverName string
	Message      string
	Target       MessageTarget
}

func (command MessageCommand) Execute() string {
	message := fmt.Sprintf(
		"[Личное сообщение] %s: %s\n",
		command.SenderName,
		command.Message,
	)

	err := command.Target.Send(message)
	if err != nil {
		return fmt.Sprintf("Unable to send private message: %s\n", err)
	}

	return fmt.Sprintf("Message sent to %s\n", command.ReceiverName)
}
