package main

import (
	"Chat/internal/server"
	"log"
	"time"
)

const autosaveInterval = time.Second

func autosave(hub *server.Hub, filename string, stop <-chan struct{}, done chan<- struct{}) {
	defer close(done)

	ticker := time.NewTicker(autosaveInterval)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return

		case <-ticker.C:
			err := server.SaveHistory(filename, hub.ExportHistory())
			if err != nil {
				log.Printf("Unable to autosave history: %s", err)
			}
		}
	}
}
