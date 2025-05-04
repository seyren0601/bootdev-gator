package main

import (
	"context"
	"database/sql"
	"gator/internal/database"
	"log"
	"strings"
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

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
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
		_, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			CreatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: *item.PublishedParsed,
			FeedID:      nextFeed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "unique") {
				// ignore fetched post
			} else {
				log.Println(err)
			}
		}
	}
	log.Println("Saved posts into database.")

	return nil
}
