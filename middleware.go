package main

import (
	"context"
	"gator/internal/database"
	"log"
	"os"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.config.Current_user_name)
		if err != nil {
			return err
		}

		err = handler(s, cmd, user)
		if err != nil {
			return err
		}

		return nil
	}
}

func middlewareLogging(handler func(s *state, cmd command) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return err
		}
		defer f.Close()

		logger := log.New(f, "[Logging] ", log.Ldate|log.Ltime)
		logger.Printf("User [%s] request command [%s] with parameters: %v\n", s.config.Current_user_name, cmd.name, cmd.parameters)

		err = handler(s, cmd)
		if err != nil {
			return err
		}

		return nil
	}
}
