package main

import (
	"context"
	"database/sql"
	"fmt"
	"gator/internal/database"
	"log"
	"time"

	"github.com/mmcdole/gofeed"
)

func fetchFeed(ctx context.Context, feedUrl string) (*gofeed.Feed, error) {
	parser := gofeed.NewParser()
	feed, err := parser.ParseURL(feedUrl)
	if err != nil {
		return nil, err
	}

	return feed, nil
}

func scrapeFeeds(s *state, user database.User) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background(), user.ID)
	if err != nil {
		return err
	}

	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		ID: nextFeed.ID,
		LastFetchedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})
	if err != nil {
		return err
	}

	feed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return err
	}

	log.Printf("Fetched feed from: %s\n", feed.Title)

	for _, item := range feed.Items {
		fmt.Printf("\t* %s\n", item.Title)
	}

	return nil
}
