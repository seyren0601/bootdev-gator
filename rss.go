package main

import (
	"context"
	"fmt"
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

func PrintFeedsItems(feeds []struct {
	Source  string
	Authors []*gofeed.Person
	Item    gofeed.Item
}) {
	fmt.Print("\033[H\033[2J")
	for _, feed := range feeds {
		fmt.Printf("Source: %s\n", feed.Source)
		fmt.Printf("Author(s): ")
		for _, author := range feed.Authors {
			fmt.Printf("%s ", author.Name)
		}
		if strings.Contains(feed.Source, "reddit") {
			fmt.Printf("\nSubreddit: %s", feed.Item.Categories[0])
		}
		fmt.Printf("\nPublished at: %s\n"+
			"Title: %s\n"+
			"Url: %s\n\n",
			feed.Item.PublishedParsed.UTC().Local().Format(time.RFC1123Z), feed.Item.Title, feed.Item.Link)
	}
}
