package main

import (
	"fmt"
	"os"
	"gomoku/internal/cli"
)

func main() {
	// Create CLI handler
	handler, err := cli.NewHandler()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing application: %v\n", err)
		os.Exit(1)
	}

	// Get command line arguments (skip program name)
	args := os.Args[1:]

	// Handle the command
	err = handler.HandleCommand(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}