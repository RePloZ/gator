package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/RePloZ/gator/internal/database"
	"github.com/RePloZ/gator/internal/handlers"
	"github.com/RePloZ/gator/internal/middleware"
	"github.com/RePloZ/gator/pkg/utils"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := utils.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	app := utils.State{Config: cfg, DB: dbQueries}

	list_commands := utils.Commands{}
	list_commands.Register("login", handlers.HandlerLogin)
	list_commands.Register("register", handlers.HandleRegister)
	list_commands.Register("reset", handlers.HandleReset)
	list_commands.Register("users", handlers.HandleAllUsers)
	list_commands.Register("agg", handlers.HandleAggregate)
	list_commands.Register("addfeed", middleware.MiddlewareLoggedIn(handlers.HandleAddFeed))
	list_commands.Register("feeds", handlers.HandleFeeds)
	list_commands.Register("follow", middleware.MiddlewareLoggedIn(handlers.HandleFollow))
	list_commands.Register("following", middleware.MiddlewareLoggedIn(handlers.HandleFollowing))
	list_commands.Register("unfollow", middleware.MiddlewareLoggedIn(handlers.HandleUnfollow))
	list_commands.Register("browse", handlers.HandleBrowse)

	args := os.Args
	if len(args) < 2 {
		log.Fatal("Cannot execute any command ! There is less than two arguments")
	}
	cmd := utils.Command{args[1], args[2:]}
	err = list_commands.Run(&app, cmd)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
