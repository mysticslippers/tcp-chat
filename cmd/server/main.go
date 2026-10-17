package main

import (
	"Chat/internal/server"
	"errors"
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

	stop := make(chan struct{})
	done := make(chan struct{})
	shutDownRequested := make(chan struct{})

	go func(listener net.Listener, stop <-chan struct{}, done chan struct{}) {
		defer close(done)

		for {
			connection, err := listener.Accept()
			if err != nil {
				if !errors.Is(err, net.ErrClosed) {
					log.Printf("Unable to accept connection: %s", err)
				}
				return
			}

			err = server.Handle(connection, stop)
			if err != nil && !errors.Is(err, net.ErrClosed) {
				log.Printf("Unable to handle connection: %s", err)
			}
		}
	}(listener, stop, done)

	go func(requested chan struct{}) {
		shutdown()
		close(requested)
	}(shutDownRequested)

	log.Println("Server launched on", address)

	select {
	case <-shutDownRequested:
	case <-done:
	}

	log.Println("Stopping server...")
	close(stop)

	if err := listener.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
		log.Printf("Unable to close listener: %s", err)
	}

	<-done
	log.Println("Server is shutdown")
}
