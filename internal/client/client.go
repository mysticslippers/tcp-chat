package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
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

	return nil
}
