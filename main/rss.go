package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/gitc-azz/bootdotdev-go-build-a-blog-aggregator/internal/database"
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
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "gator")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rssFeed RSSFeed
	err = xml.Unmarshal(bytes, &rssFeed)
	if err != nil {
		return nil, err
	}

	unescapeString(&rssFeed)

	return &rssFeed, nil
}

func unescapeString(feed *RSSFeed) {
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for idx, _ := range feed.Channel.Item {
		feed.Channel.Item[idx].Title = html.UnescapeString(feed.Channel.Item[idx].Title)
		feed.Channel.Item[idx].Description = html.UnescapeString(feed.Channel.Item[idx].Description)
	}
}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("Failed to get the next feed to fetch, err= %v", err)
	}

	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		ID:            feed.ID,
		LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to mark this feed: %v at %v, as fetched, err= %v", feed.Name, feed.Url, err)
	}

	rssFeed, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("failed to fetchFeed at %v, err= %v", feed.Url, err)
	}

	if len(rssFeed.Channel.Item) == 0 {
		fmt.Println("RSS feed has zero Item")
	} else {
		fmt.Println("RSS feed -", feed.Name, "- has the following Item's title:")
	}
	for _, item := range rssFeed.Channel.Item {
		fmt.Println(" *", item.Title)
	}

	return nil
}
