package middleware

import (
	"context"

	"github.com/RePloZ/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *models.state, cmd models.command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		ctx := context.Background()
		defer ctx.Done()

		usr, err := s.db.GetUser(ctx, s.config.CurrentUserName)
		if err != nil {
			return err
		}

		return handler(s, cmd, usr)
	}
}
