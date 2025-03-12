package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/willthefoollearn/gator/internal/database"

	_ "github.com/lib/pq"
	"github.com/willthefoollearn/gator/internal/config"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	con, err := config.Read()
	if err != nil {
		log.Fatalf("Unable to read config file: %v\n", err)
	}

	newState := &state{
		cfg: &con,
	}

	db, err := sql.Open("postgres", newState.cfg.DbURL)
	if err != nil {
		log.Fatalf("Unable to open database: %v\n", err)
	}
	defer db.Close()

	dbQueries := database.New(db)
	newState.db = dbQueries
	newCommands := commands{
		callback: make(map[string]func(*state, command) error),
	}
	commandArgs := os.Args[1:]

	if len(commandArgs) < 1 {
		log.Fatal("Not enough arguments\n")
	}

	newCommands.register("login", handlerLogin)
	newCommands.register("register", handlerRegister)
	newCommands.register("reset", handlerReset)
	newCommands.register("users", handlerList)
	newCommands.register("agg", handlerAgg)
	newCommands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	newCommands.register("feeds", handlerFeeds)
	newCommands.register("follow", middlewareLoggedIn(handlerFollow))
	newCommands.register("following", middlewareLoggedIn(handlerFollowing))
	newCommands.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	newCommands.register("browse", middlewareLoggedIn(handlerBrowse))

	cmd := command{
		name: commandArgs[0],
		args: commandArgs[1:],
	}

	err = newCommands.run(newState, cmd)
	if err != nil {
		log.Fatalf("Command doesn't exist: %v\n", err)
	}

	con, err = config.Read()
	if err != nil {
		log.Fatalf("Unable to read config file: %v\n", err)
	}
}
