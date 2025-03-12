package main

import "errors"

type command struct {
	name string
	args []string
}

type commands struct {
	name     string
	callback map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.name = name
	c.callback[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.callback[cmd.name]
	if !ok {
		return errors.New("command not found")
	}
	
	return f(s, cmd)
}
