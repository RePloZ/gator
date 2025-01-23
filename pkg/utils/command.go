package utils

import (
	"fmt"

	"github.com/RePloZ/gator/internal/database"
)

type Command struct {
	Name string
	Args []string
}

type State struct {
	DB     *database.Queries
	Config Config
}

type Commands struct {
	registeredCommands map[string]func(*State, Command) error
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	if c.registeredCommands == nil {
		c.registeredCommands = make(map[string]func(*State, Command) error)
	}
	c.registeredCommands[name] = f
}

func (c *Commands) Run(s *State, cmd Command) error {
	handler, exists := c.registeredCommands[cmd.Name]
	if !exists {
		return fmt.Errorf("command not found: %s", cmd.Name)
	}
	return handler(s, cmd)
}
