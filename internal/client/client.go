package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func Client(address string) error {
	connection, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	defer func(connection net.Conn) {
		err := connection.Close()
		if err != nil {
			log.Printf("Error closing connection: %s", err)
		}
	}(connection)

	reader := bufio.NewReader(connection)

	message, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	fmt.Print(message)

	fmt.Print("Please, input your message: ")
	inputReader := bufio.NewReader(os.Stdin)
	input, err := inputReader.ReadString('\n')
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(connection, input)
	if err != nil {
		return err
	}

	reply, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	fmt.Print("Reply from server: ", reply)

	return nil
}
