package main

import (
	"Chat/internal/client"
	"log"
)

func main() {
	const address = "localhost:9000"

	err := client.Client(address)
	if err != nil {
		log.Fatalf("Client error: %s", err)
	}
}
