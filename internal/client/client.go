package client

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
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

	inputReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Please, input your message: ")
		input, err := inputReader.ReadString('\n')
		if strings.TrimSpace(input) == "quit" {
			return nil
		}
		if err == io.EOF {
			return nil
		}
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
	}
}
