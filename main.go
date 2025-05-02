package main

import (
	"database/sql"
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	// Read args from command-line call
	// If invalid number of args, exit with code 1
	args := os.Args
	if len(args) < 2 {
		fmt.Println("command required")
		os.Exit(1)
	}

	// Read config file into Config struct
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	// Connect to database
	db, err := sql.Open("postgres", cfg.Db_url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Create Queries (generated database package) instance
	dbQueries := database.New(db)

	// Create state instance
	// and store the Queries instance into state
	st := state{
		config: &cfg,
		db:     dbQueries,
	}

	// Create commands instance and register handlers
	commands := NewCommands().
		register("login", handlerLogin).
		register("register", handlerRegister).
		register("reset", handlerReset)

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
