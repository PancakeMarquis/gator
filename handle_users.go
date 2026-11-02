package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/pancakemarquis/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("Must include username")
	}
	name := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("couldn't find user: %w", err)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return err
	}
	fmt.Println("User has been set")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("usage: %v <name>", cmd.name)
	}
	userName := cmd.args[0]
	id := uuid.New()

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{ID: id, CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(), Name: userName})
	if err != nil {
		return fmt.Errorf("couldn't create user: %w", err)
	}

	err = s.cfg.SetUser(userName)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}
	fmt.Println("User created successfully:")
	printUser(user)
	return nil
}

func printUser(user database.User) {
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)
}
