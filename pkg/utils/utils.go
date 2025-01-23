package utils

import (
	"fmt"

	"github.com/RePloZ/gator/internal/database"
)

type command struct {
	Name string
	Args []string
}

type state struct {
	db     *database.Queries
	config utils.Config
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	if c.registeredCommands == nil {
		c.registeredCommands = make(map[string]func(*state, command) error)
	}
	c.registeredCommands[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	handler, exists := c.registeredCommands[cmd.Name]
	if !exists {
		return fmt.Errorf("command not found: %s", cmd.Name)
	}
	return handler(s, cmd)
}
