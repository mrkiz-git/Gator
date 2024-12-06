package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
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

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	if feedURL == "" {
		err := fmt.Errorf("missing URL")
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)

	if err != nil {
		log.Printf("failed creating http request to url %s", feedURL)
		return nil, err
	}

	req.Header.Add("User-Agent", "gator")
	httpClient := &http.Client{Timeout: 10 * time.Second}
	resp, err := httpClient.Do(req)

	if err != nil {
		log.Printf("faild http request")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status %d for URL %s", resp.StatusCode, feedURL)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {

		return nil, fmt.Errorf("Error reading response: %w", err)

	}

	rssFeedOuput := new(RSSFeed)
	err = xml.Unmarshal(body, rssFeedOuput)
	if err != nil {
		return nil, fmt.Errorf("Error reading decoding XML: %w", err)
	}

	return rssFeedOuput, nil

}
