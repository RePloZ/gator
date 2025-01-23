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
	defer ctx.Done()
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

func handleReset(s *state, _ command) error {
	ctx := context.Background()
	defer ctx.Done()
	err := s.db.DeleteAllUsers(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("All users was deleted !")
	os.Exit(0)
	return nil
}

func handleAllUsers(s *state, _ command) error {
	ctx := context.Background()
	defer ctx.Done()
	allUsers, err := s.db.GetUsers(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("Get all users !\n")
	for _, user := range allUsers {
		fmt.Printf("* %s", user.Name)
		if user.Name == s.config.CurrentUserName {
			fmt.Printf(" (current)")
		}
		fmt.Print("\n")
	}
	os.Exit(0)
	return nil
}

func handleUuid(s *state, _ command) error {
	fmt.Println(uuid.New())
	return nil
}
