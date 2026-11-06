package main

import (
	"Chat/internal/server"
	"errors"
	"log"
	"net"
	"os"
	"sync"
)

func main() {
	const (
		filename = "history.json"
		address  = "localhost:9000"
	)

	hub := server.NewHub()

	history, err := server.LoadHistory(filename)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Fatalf("Unable to load history: %s", err)
		}
	} else {
		err = hub.RestoreHistory(history)
		if err != nil {
			log.Fatalf("Unable to restore history: %s", err)
		}
	}

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Unable to start server: %s", err)
	}

	var handlers sync.WaitGroup
	stop := make(chan struct{})
	done := make(chan struct{})
	autosaveDone := make(chan struct{})
	shutDownRequested := make(chan struct{})

	go autosave(hub, filename, stop, autosaveDone)

	go func(listener net.Listener, stop <-chan struct{}, done chan struct{}, hub *server.Hub) {
		defer close(done)

		for {
			connection, err := listener.Accept()
			if err != nil {
				if !errors.Is(err, net.ErrClosed) {
					log.Printf("Unable to accept connection: %s", err)
				}
				return
			}

			handlers.Add(1)

			go func(connection net.Conn, stop <-chan struct{}, hub *server.Hub) {
				defer handlers.Done()

				err := server.Handle(connection, stop, hub)
				if err != nil && !errors.Is(err, net.ErrClosed) {
					log.Printf("Unable to handle connection: %s", err)
				}
			}(connection, stop, hub)
		}
	}(listener, stop, done, hub)

	go shutdown(shutDownRequested)

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
	handlers.Wait()
	<-autosaveDone

	err = server.SaveHistory(filename, hub.ExportHistory())
	if err != nil {
		log.Printf("Unable to save history: %s", err)
	}

	log.Println("Server is shutdown")
}
