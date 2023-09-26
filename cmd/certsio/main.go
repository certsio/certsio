package main

import (
	"log"

	"github.com/certsio/certsio/internal/runner/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("error executing command: %v", err)
	}
}
