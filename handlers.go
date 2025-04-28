package main

import (
	"errors"
	"gator/internal/config"
)

type state struct {
	config *config.Config
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.parameters) != 1 {
		return errors.New("login command expects 1 parameter")
	}

	username := cmd.parameters[0]
	s.config.Current_user_name = username
	err := s.config.SetUser()
	if err != nil {
		return err
	}

	return nil
}
