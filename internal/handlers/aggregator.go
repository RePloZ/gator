package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/RePloZ/gator/internal/database"
	"github.com/RePloZ/gator/pkg/utils"
	"github.com/google/uuid"
)

func HandleAddFeed(s *utils.State, cmd utils.Command, usr database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("no Argument")
	}

	ctx := context.Background()
	defer ctx.Done()

	createFeeds := database.CreateFeedsParams{
		ID:        uuid.New(),
		Name:      sql.NullString{String: cmd.Args[0], Valid: true},
		Url:       sql.NullString{String: cmd.Args[1], Valid: true},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    usr.ID,
	}
	feed, err := s.DB.CreateFeeds(ctx, createFeeds)
	if err != nil {
		return err
	}

	createFeedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    usr.ID,
		FeedID:    feed.ID,
	}
	_, err = s.DB.CreateFeedFollow(ctx, createFeedFollow)
	if err != nil {
		return fmt.Errorf("cannot create feed follow : %w", err)
	}

	fmt.Println(feed)
	return nil
}

func HandleAggregate(s *utils.State, cmd utils.Command) error {
	if len(cmd.Args) < 1 || len(cmd.Args) > 2 {
		return fmt.Errorf("usage: %v <time_between_reqs>", cmd.Name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w", err)
	}

	log.Printf("Collecting feeds every %s...", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		utils.ScrapeFeeds(s)
	}
}

func HandleFeeds(s *utils.State, cmd utils.Command) error {
	ctx := context.Background()
	defer ctx.Done()
	feeds, err := s.DB.GetFeedsInformation(ctx)
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Println(feed)
	}
	return nil
}

func HandleFollow(s *utils.State, cmd utils.Command, usr database.User) error {
	ctx := context.Background()
	defer ctx.Done()

	url := cmd.Args[0]
	feed, err := s.DB.GetFeeedByUrl(ctx, sql.NullString{String: url, Valid: true})
	if err != nil {
		return err
	}

	createFeedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    usr.ID,
		FeedID:    feed.ID,
	}
	feed_follow, err := s.DB.CreateFeedFollow(ctx, createFeedFollow)
	if err != nil {
		return err
	}

	fmt.Println(feed_follow)
	return nil
}

func HandleFollowing(s *utils.State, cmd utils.Command, usr database.User) error {
	ctx := context.Background()
	defer ctx.Done()

	usr, err := s.DB.GetUser(ctx, s.Config.CurrentUserName)
	if err != nil {
		return err
	}

	feed_follow, err := s.DB.GetFeedFollowsForUser(ctx, usr.ID)
	if err != nil {
		return err
	}

	fmt.Println(feed_follow)
	return nil
}

func HandleUnfollow(s *utils.State, cmd utils.Command, usr database.User) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("no Argument")
	}
	ctx := context.Background()
	defer ctx.Done()

	url := cmd.Args[0]
	feed, err := s.DB.GetFeeedByUrl(ctx, sql.NullString{Valid: true, String: url})
	if err != nil {
		return err
	}

	deleteFeedFollowsForUserParams := database.DeleteFeedFollowsForUserParams{
		UserID: usr.ID,
		FeedID: feed.ID,
	}
	err = s.DB.DeleteFeedFollowsForUser(ctx, deleteFeedFollowsForUserParams)

	if err != nil {
		return err
	}

	return nil
}

func HandleBrowse(s *utils.State, cmd utils.Command) error {
	log.Println(len(cmd.Args))
	if len(cmd.Args) > 2 {
		return fmt.Errorf("Too much arguments")
	}

	limit, err := strconv.ParseInt(cmd.Args[0], 10, 64)
	if err != nil {
		limit = 2
	}

	posts, err := s.DB.GetPosts(context.Background(), int32(limit))
	if err != nil {
		return err
	}

	for _, post := range posts {
		log.Println(post)
	}

	return nil
}
