package main

import (
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
			log.Fatalf("Unable to close listener: %s", err)
		}
	}(listener)

	log.Println("Server started on", address)
	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Printf("Unable to accept connection: %s", err)
			return
		}
		log.Printf("Accepted connection from %s", connection.RemoteAddr())

		closeErr := connection.Close()
		if closeErr != nil {
			log.Printf("Unable to close connection: %s", closeErr)
		}
	}
}
