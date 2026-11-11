package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/pancakemarquis/gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {

	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <time_between_reqs>", cmd.name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("Unable to convert time: %w", err)
	}
	log.Printf("Collecting feeds every %s...", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapFeeds(s)
	}

}

func scrapFeeds(s *state) {
	nextFeed, err := s.db.GetNextFeedtoFetch(context.Background())
	if err != nil {
		log.Printf("Couldn't mark feed %s fetched: %v", nextFeed.Name, err)
		return
	}
	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{LastFetchedAt: sql.NullTime{Time: time.Now().UTC()}, UpdatedAt: time.Now().UTC(), ID: nextFeed.ID})

	if err != nil {
		return
	}

	feeds, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		log.Printf("Couldn't collect feed %s: %v", nextFeed.Name, err)
		return
	}

	for _, item := range feeds.Channel.Item {
		fmt.Printf("Found post: %s\n", item.Title)
	}
	log.Printf("Feed %s collected, %v posts found", nextFeed.Name, len(feeds.Channel.Item))

}
