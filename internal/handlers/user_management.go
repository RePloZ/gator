package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/RePloZ/gator/internal/database"
	"github.com/RePloZ/gator/pkg/utils"
	"github.com/google/uuid"
)

func HandlerLogin(s *utils.State, cmd utils.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("no Argument")
	}

	if cmd.Args[0] == "unknown" {
		return fmt.Errorf("user not found")
	}

	err := s.Config.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Printf("User has been set !")
	return nil
}

func HandleRegister(s *utils.State, cmd utils.Command) error {
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
	_, err := s.DB.CreateUser(ctx, user)
	if err != nil {
		return err
	}
	fmt.Printf("User has been created !")
	s.Config.SetUser(user.Name)
	return nil
}

func HandleReset(s *utils.State, _ utils.Command) error {
	ctx := context.Background()
	defer ctx.Done()
	err := s.DB.DeleteAllUsers(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("All users was deleted !")
	return nil
}

func HandleAllUsers(s *utils.State, _ utils.Command) error {
	ctx := context.Background()
	defer ctx.Done()
	allUsers, err := s.DB.GetUsers(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("Get all users !\n")
	for _, user := range allUsers {
		fmt.Printf("* %s", user.Name)
		if user.Name == s.Config.CurrentUserName {
			fmt.Printf(" (current)")
		}
		fmt.Print("\n")
	}
	return nil
}
