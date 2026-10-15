package main

import (
	"Chat/internal/client"
	"log"
)

func main() {
	const address = "127.0.0.1:9000"

	err := client.Client(address)
	if err != nil {
		log.Fatalf("Client error: %s", err)
	}
}
