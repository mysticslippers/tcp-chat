package main

import (
	"Chat/internal/server"
	"log"
	"net"
)

func main() {
	const address = "127.0.0.1:9000"

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Unable to start server: %s", err)
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Printf("Unable to close listener: %s", err)
		}
	}(listener)

	log.Println("Server started on", address)
	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Printf("Unable to accept connection: %s", err)
			return
		}

		err = server.Handle(connection)
		if err != nil {
			log.Printf("Unable to handle connection: %s", err)
		}
	}
}
