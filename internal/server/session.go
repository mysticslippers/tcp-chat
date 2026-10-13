package server

import (
	"fmt"
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
	_, err := fmt.Fprintln(connection, "Hello!")
	if err != nil {
		return err
	}

	return nil
}
