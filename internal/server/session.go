package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

func Handle(connection net.Conn) error {
	defer func(connection net.Conn) {
		err := connection.Close()
		if err != nil {
			log.Printf("Error closing connection: %s", err)
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

		_, err = fmt.Fprintln(connection, message)
		if err != nil {
			return err
		}
	}
}
