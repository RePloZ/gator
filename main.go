package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/RePloZ/gator/internal/config"
	"github.com/RePloZ/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	app := cli.state{config: cfg, db: dbQueries}

	list_commands := commands{}
	list_commands.register("login", handlers.handlerLogin)
	list_commands.register("register", handlers.handleRegister)
	list_commands.register("reset", handlers.handleReset)
	list_commands.register("users", handlers.handleAllUsers)
	list_commands.register("agg", handlers.handleAggregate)
	list_commands.register("addfeed", middleware.middlewareLoggedIn(handlers.handleAddFeed))
	list_commands.register("feeds", handlers.handleFeeds)
	list_commands.register("follow", middleware.middlewareLoggedIn(handlers.handleFollow))
	list_commands.register("following", middleware.middlewareLoggedIn(handlers.handleFollowing))
	list_commands.register("unfollow", middleware.middlewareLoggedIn(handlers.handleUnfollow))
	list_commands.register("browse", handlers.handleBrowse)

	args := os.Args
	if len(args) < 2 {
		log.Fatal("Cannot execute any command ! There is less than two arguments")
	}
	cmd := command{args[1], args[2:]}
	err = list_commands.run(&app, cmd)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
