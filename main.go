package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/pancakemarquis/gator/internal/config"
	"github.com/pancakemarquis/gator/internal/database"

	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {

	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatalf("error connecting to Database: %v", err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	s := state{db: dbQueries, cfg: &cfg}
	cmds := commands{
		registeredCommands: map[string]func(*state, command) error{},
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerListUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerListFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFeedFollow))
	cmds.register("following", middlewareLoggedIn(handlerListFeedFollow))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollowFeedFollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	s.db = dbQueries

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
