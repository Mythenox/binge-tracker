package app

import (
	"errors"
)

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	RegisteredCmds map[string]func(*State, Command) error
}

func (c *Commands) Run(s *State, cmd Command) error {
	f, ok := c.RegisteredCmds[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}

	return f(s, cmd)
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.RegisteredCmds[name] = f
}
