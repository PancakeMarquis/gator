package main

import (
	"errors"
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("Must include username")
	}
	err := s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	}
	fmt.Println("User has been set")
	return nil
}
