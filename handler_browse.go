package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/pancakemarquis/gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int

	if len(cmd.args) == 1 {
		newLimit, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("Unable to convert to int: %w", err)
		}
		if newLimit <= 0 {
			return fmt.Errorf("Limit should not be zero or negative")
		}
		limit = newLimit
	} else if len(cmd.args) > 1 {
		return fmt.Errorf("usage: %s Optional:<limit>", cmd.name)
	} else {
		limit = 2
	}

	posts, err := s.db.GetPostsByUser(context.Background(), database.GetPostsByUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("Unable to get posts from user: %w", err)
	}
	for _, post := range posts {
		printPosts(post)

	}

	return nil

}

func printPosts(post database.GetPostsByUserRow) {
	fmt.Printf("%s from %s\n", post.PublishedAt.Time.Format("Mon Jan 2"), post.FeedName)
	fmt.Printf("--- %s ---\n", post.Title)
	fmt.Printf("    %v\n", post.Description.String)
	fmt.Printf("Link: %s\n", post.Url)
	fmt.Println("=====================================")
}
