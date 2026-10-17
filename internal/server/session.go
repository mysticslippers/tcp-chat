package server

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

func Handle(connection net.Conn, stop <-chan struct{}) error {
	defer func(connection net.Conn) {
		err := connection.Close()
		if err != nil && !errors.Is(err, net.ErrClosed) {
			log.Printf("Error closing connection: %s", err)
		}
	}(connection)

	done := make(chan struct{})
	defer close(done)

	go func(connection net.Conn) {
		select {
		case <-stop:
			err := connection.Close()
			if err != nil && !errors.Is(err, net.ErrClosed) {
				log.Printf("Error closing connection: %s", err)
			}

		case <-done:
		}
	}(connection)

	log.Printf("Accepted connection from %s", connection.RemoteAddr())
	_, err := fmt.Fprintln(connection, "You have successfully connected to the server")
	if err != nil {
		return err
	}

	reader := bufio.NewReader(connection)
	for {
		message, err := reader.ReadString('\n')
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		_, err = fmt.Fprint(connection, message)
		if err != nil {
			return err
		}
	}
}
