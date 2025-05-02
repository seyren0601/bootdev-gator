package main

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
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

func fetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {
	client := http.Client{}

	req, err := http.NewRequest("GET", feedUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator/0.4")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var feed RSSFeed
	err = xml.Unmarshal(bytes, &feed)
	if err != nil {
		return nil, err
	}

	// Decode escape characters (code => char) using html.UnescapeString
	// Opposite is html.EscapeString (char => code)
	//
	// Decode texts in channel titles and descriptions
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	// Decode texts in items titles and descriptions
	for i, item := range feed.Channel.Item {
		feed.Channel.Item[i].Description = html.UnescapeString(item.Description)
		feed.Channel.Item[i].Title = html.UnescapeString(item.Title)
	}
	return &feed, nil
}
