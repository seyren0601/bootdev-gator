package main

import (
	"context"
	"errors"
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
	"time"
)

type state struct {
	config *config.Config
	db     *database.Queries
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.parameters) != 1 {
		return errors.New("login command expects 1 parameter")
	}

	username := cmd.parameters[0]

	user, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return errors.New("username doesn't exist")
	}

	s.config.Current_user_name = user.Name
	err = s.config.SetUser()
	if err != nil {
		return err
	}

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.parameters) != 1 {
		return errors.New("register command expects 1 parameters")
	}

	username := cmd.parameters[0]

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	})

	if err != nil {
		return err
	}

	fmt.Printf(`User created successfully
	UUID: %s
	CreatedAt: %s
	UpdatedAt: %s
	Name: %s\n
`, user.ID, user.CreatedAt.Local().String(), user.UpdatedAt.Local().String(), user.Name)

	s.config.Current_user_name = user.Name
	s.config.SetUser()

	return nil
}
