package main

import (
	"fmt"
	"gator/internal/config"
	"os"
)

func main() {
	// Read args from command-line call
	// If invalid number of args, exit with code 1
	args := os.Args
	if len(args) < 2 {
		fmt.Println("command required")
		os.Exit(1)
	}

	// Read config file
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	// Create state instance
	st := state{
		config: &cfg,
	}

	// Create commands instance and register handlers
	commands := commands{
		callbacks: map[string]func(*state, command) error{},
	}
	commands.register("login", handlerLogin)

	// Create command instance based on args
	cmd := command{
		name:       args[1],
		parameters: args[2:],
	}

	// Run command
	err = commands.run(&st, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Exit code 0 (correct)
	os.Exit(0)
}
