package middleware

import (
	"context"

	"github.com/RePloZ/gator/internal/database"
	"github.com/RePloZ/gator/pkg/utils"
)

func MiddlewareLoggedIn(handler func(s *utils.State, cmd utils.Command, user database.User) error) func(*utils.State, utils.Command) error {
	return func(s *utils.State, cmd utils.Command) error {
		ctx := context.Background()
		defer ctx.Done()

		usr, err := s.DB.GetUser(ctx, s.Config.CurrentUserName)
		if err != nil {
			return err
		}

		return handler(s, cmd, usr)
	}
}
