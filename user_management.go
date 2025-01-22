package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/RePloZ/gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("no Argument")
	}

	if cmd.Args[0] == "unknown" {
		return fmt.Errorf("user not found")
	}

	err := s.config.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Printf("User has been set !")
	os.Exit(0)
	return nil
}

func handleRegister(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("no Argument")
	}

	ctx := context.Background()
	user := database.CreateUserParams{
		ID:        uuid.New(),
		Name:      cmd.Args[0],
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err := s.db.CreateUser(ctx, user)
	if err != nil {
		return err
	}
	fmt.Printf("User has been created !")
	s.config.SetUser(user.Name)
	os.Exit(0)
	return nil
}
