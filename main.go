package main

import (
	"blog-aggregator/internal/config"
	"blog-aggregator/internal/database"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

const dbURL = "postgres://postgres:postgres@localhost:5432/gator"

func main() {
	db, err := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)

	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config: ", err)
		os.Exit(1)
	}

	st := state{cfg: &cfg, db: dbQueries}
	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Not enough arguments")
		os.Exit(1)
	}
	commandName := args[1]
	args = args[2:]
	err = cmds.run(&st, command{name: commandName, args: args})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
