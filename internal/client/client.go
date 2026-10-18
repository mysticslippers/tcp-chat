package client

import (
	"bufio"
	"errors"
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
	fmt.Println("Enter messages. Type quit to disconnect")

	inputDone := make(chan error, 1)
	serverDone := make(chan error, 1)

	go func(connection net.Conn, inputDone chan<- error) {
		inputDone <- send(connection)
	}(connection, inputDone)

	go func(reader *bufio.Reader, serverDone chan<- error) {
		for {
			reply, err := reader.ReadString('\n')
			if err != nil {
				serverDone <- err
				return
			}

			fmt.Print("Reply from server: ", reply)
			fmt.Print("> ")
		}
	}(reader, serverDone)

	select {
	case err := <-inputDone:
		return err

	case err := <-serverDone:
		if errors.Is(err, io.EOF) {
			fmt.Println("Server closed the connection")
			return nil
		}

		return fmt.Errorf("connection to server lost: %w", err)
	}
}

func send(connection net.Conn) error {
	inputReader := bufio.NewReader(os.Stdin)
	for {
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
	}
}
