package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pancakemarquis/gator/internal/database"
)

func handlerFeedFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.name)
	}

	url := cmd.args[0]
	feed, err := s.db.GetFeedsByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("Couldn't get feed: %w", err)
	}

	feedFollows, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Couldn't create feedFollow record: %w", err)
	}
	fmt.Println("Feed follow created:")
	printFeedFollow(feedFollows.UserName, feedFollows.FeedName)

	return nil

}
func handlerListFeedFollow(s *state, cmd command, user database.User) error {
	listFollow, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("Couldn't get feed_follow: %w", err)
	}
	if len(listFollow) == 0 {
		fmt.Println("No feed follows found for this user.")
		return nil
	}
	fmt.Printf("Feed follows for user %s:\n", user.Name)
	for _, feed := range listFollow {
		fmt.Printf("Feed name: %s\n", feed.FeedName)
	}

	return nil

}
func handlerUnfollowFeedFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <feed_url>", cmd.name)
	}

	url := cmd.args[0]

	feed, err := s.db.GetFeedsByURL(context.Background(), url)

	if err != nil {
		return fmt.Errorf("Couldn't get feed: %w", err)
	}

	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{FeedID: feed.ID, UserID: user.ID})
	if err != nil {
		return fmt.Errorf("Couldn't delete feed follow: %w", err)
	}
	fmt.Printf("Feed %s has been unfollowed by User: %s", feed.Name, user.Name)

	return nil

}

func printFeedFollow(username, feedname string) {
	fmt.Printf("* User:          %s\n", username)
	fmt.Printf("* Feed:          %s\n", feedname)
}
