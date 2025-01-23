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

func handleAddFeed(s *state, cmd command) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("no Argument")
	}

	ctx := context.Background()
	defer ctx.Done()
	user, err := s.db.GetUser(ctx, s.config.CurrentUserName)
	if err != nil {
		return err
	}

	createFeeds := database.CreateFeedsParams{
		ID:        uuid.New(),
		Name:      sql.NullString{String: cmd.Args[0], Valid: true},
		Url:       sql.NullString{String: cmd.Args[1], Valid: true},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
	}
	feed, err := s.db.CreateFeeds(ctx, createFeeds)
	if err != nil {
		return err
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
	fmt.Println(rss)
	os.Exit(0)
	return nil
}
