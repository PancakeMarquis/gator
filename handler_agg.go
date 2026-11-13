package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
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
		publishedAt := sql.NullTime{}
		pubDate, err := time.Parse(time.RFC1123Z, fmt.Sprint(item.PubDate))
		if err != nil {
			log.Printf("Couldn't parse published date  %v", err)
		} else {
			publishedAt.Time = pubDate
			publishedAt.Valid = true
		}
		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: true},
			PublishedAt: sql.NullTime{Time: publishedAt.Time, Valid: publishedAt.Valid},
			FeedID:      nextFeed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				log.Printf("Post with Url already exists")
				continue
			} else {
				log.Printf("Couldn't create post: %v", err)
				return
			}
		}
	}

}
