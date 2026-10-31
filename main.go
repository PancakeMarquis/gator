package main

import (
	"log"
	"os"

	"github.com/pancakemarquis/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	s := state{cfg: &cfg}
	cmds := commands{
		list: map[string]func(*state, command) error{},
	}
	cmds.register("login", handlerLogin)

	input := os.Args
	if len(input) < 2 {
		log.Fatal("No arguments passed")
	}
	cmdName := input[1]
	params := input[2:]
	cmd := command{
		name: cmdName,
		args: params,
	}

	err = cmds.run(&s, cmd)
	if err != nil {
		log.Fatal(err)
	}

}
