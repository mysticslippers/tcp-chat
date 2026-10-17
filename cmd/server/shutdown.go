package main

import (
	"bufio"
	"log"
	"os"
	"os/signal"
	"strings"
)

func shutdown(requested chan struct{}) {
	command := make(chan struct{})
	interrupt := make(chan os.Signal, 1)

	signal.Notify(interrupt, os.Interrupt)
	defer signal.Stop(interrupt)

	go func(command chan struct{}) {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			if strings.TrimSpace(scanner.Text()) == "exit" {
				close(command)
				return
			}
		}
		if err := scanner.Err(); err != nil {
			log.Printf("Unable to read console: %s", err)
		}
	}(command)

	select {
	case <-command:
	case <-interrupt:
	}

	close(requested)
}
