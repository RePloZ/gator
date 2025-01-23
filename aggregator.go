package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log"
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

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "gator")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rssFeed RSSFeed
	err = xml.Unmarshal(dat, &rssFeed)
	if err != nil {
		return nil, err
	}

	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	for i, item := range rssFeed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		rssFeed.Channel.Item[i] = item
	}

	return &rssFeed, nil
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
		scrapeFeeds(s)
	}
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

func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("Couldn't get next feeds to fetch", err)
		return
	}
	log.Println("Found a feed to fetch!")
	scrapeFeed(s.db, feed)
}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Couldn't mark feed %s fetched: %v", feed.Name, err)
		return
	}

	feedData, err := fetchFeed(context.Background(), feed.Url.String)
	if err != nil {
		log.Printf("Couldn't collect feed %s: %v", feed.Name, err)
		return
	}
	for _, item := range feedData.Channel.Item {
		fmt.Printf("Found post: %s\n", item.Title)
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))
}
