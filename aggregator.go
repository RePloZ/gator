package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/RePloZ/gator/internal/database"
	"github.com/google/uuid"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (rss *RSSFeed, err error) {
	rss = &RSSFeed{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	err = xml.Unmarshal(data, rss)
	if err != nil {
		return
	}

	return
}

func handleAddFeed(s *state, cmd command, usr database.User) error {
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
	feed, err := s.db.CreateFeeds(ctx, createFeeds)
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
	_, err = s.db.CreateFeedFollow(ctx, createFeedFollow)
	if err != nil {
		return fmt.Errorf("cannot create feed follow : %w", err)
	}

	fmt.Println(feed)
	os.Exit(0)
	return nil
}

func handleAggregate(s *state, cmd command) error {
	ctx := context.Background()
	defer ctx.Done()
	rss, err := fetchFeed(ctx, "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	for _, item := range rss.Channel.Item {
		(&item).Title = html.EscapeString(item.Title)
		(&item).Description = html.EscapeString(item.Description)
	}

	duration, err := time.ParseDuration("1m")
	if err != nil {
		return err
	}
	ticker := time.NewTicker(duration)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

	os.Exit(0)
	return nil
}

func handleFeeds(s *state, cmd command) error {
	ctx := context.Background()
	defer ctx.Done()
	feeds, err := s.db.GetFeedsInformation(ctx)
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Println(feed)
	}
	return nil
}

func handleFollow(s *state, cmd command, usr database.User) error {
	ctx := context.Background()
	defer ctx.Done()

	url := cmd.Args[0]
	feed, err := s.db.GetFeeedByUrl(ctx, sql.NullString{String: url, Valid: true})
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
	feed_follow, err := s.db.CreateFeedFollow(ctx, createFeedFollow)
	if err != nil {
		return err
	}

	fmt.Println(feed_follow)
	os.Exit(0)
	return nil
}

func handleFollowing(s *state, cmd command, usr database.User) error {
	ctx := context.Background()
	defer ctx.Done()

	usr, err := s.db.GetUser(ctx, s.config.CurrentUserName)
	if err != nil {
		return err
	}

	feed_follow, err := s.db.GetFeedFollowsForUser(ctx, usr.ID)
	if err != nil {
		return err
	}

	fmt.Println(feed_follow)
	os.Exit(0)
	return nil
}

func handleUnfollow(s *state, cmd command, usr database.User) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("no Argument")
	}
	ctx := context.Background()
	defer ctx.Done()

	url := cmd.Args[0]
	feed, err := s.db.GetFeeedByUrl(ctx, sql.NullString{Valid: true, String: url})
	if err != nil {
		return err
	}

	deleteFeedFollowsForUserParams := database.DeleteFeedFollowsForUserParams{
		UserID: usr.ID,
		FeedID: feed.ID,
	}
	err = s.db.DeleteFeedFollowsForUser(ctx, deleteFeedFollowsForUserParams)

	if err != nil {
		return err
	}

	os.Exit(0)
	return nil
}

func scrapeFeeds(s *state) error {
	ctx := context.Background()
	defer ctx.Done()

	feed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		return err
	}

	err = s.db.MarkFeedFetched(
		ctx,
		database.MarkFeedFetchedParams{
			ID:        feed.ID,
			UpdatedAt: time.Now(),
		},
	)
	if err != nil {
		return err
	}

	currentFeed, err := s.db.GetFeeedByUrl(ctx, feed.Url)
	if err != nil {
		return err
	}

	fmt.Println(currentFeed.Name)
	return nil
}
