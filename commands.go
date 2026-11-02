package main

import "errors"

type command struct {
	name string
	args []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.registeredCommands[cmd.name]
	if !ok {
		return errors.New("Command does not exist")
	}
	err := handler(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}
