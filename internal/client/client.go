package client

import (
	"bufio"
	"encoding/json"
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

	decoder := json.NewDecoder(connection)
	var message string

	err = decoder.Decode(&message)
	if err != nil {
		return err
	}
	fmt.Print(message)
	fmt.Println("Enter messages. Type /quit to disconnect")
	fmt.Print("> ")

	inputDone := make(chan error, 1)
	serverDone := make(chan error, 1)

	go func(connection net.Conn, inputDone chan<- error) {
		inputDone <- send(connection)
	}(connection, inputDone)

	go func(decoder *json.Decoder, serverDone chan<- error) {
		for {
			var reply string
			err := decoder.Decode(&reply)
			if err != nil {
				serverDone <- err
				return
			}

			fmt.Print(reply)
			fmt.Print("> ")
		}
	}(decoder, serverDone)

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
		if strings.TrimSpace(input) == "/quit" {
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
