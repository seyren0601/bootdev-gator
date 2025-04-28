package main

import (
	"fmt"
)

type command struct {
	name       string
	parameters []string
}

type commands struct {
	callbacks map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	if command, found := c.callbacks[cmd.name]; !found {
		return fmt.Errorf("failed to run command [%s]. command not found", cmd.name)
	} else {
		err := command(s, cmd)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) error {
	c.callbacks[name] = f

	return nil
}
