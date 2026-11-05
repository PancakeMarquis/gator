package main

import (
	"context"
	"fmt"
)

func handlerAgg(s *state, cmd command) error {

	feedURL := "https://www.wagslane.dev/index.xml"

	rss, err := fetchFeed(context.Background(), feedURL)
	if err != nil {
		return fmt.Errorf("couldn't get feed: %w", err)
	}
	fmt.Printf("Feed: %+v\n", rss)
	return nil
}
