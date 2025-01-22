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

	app := state{config: cfg, db: dbQueries}

	list_commands := commands{}
	list_commands.register("login", handlerLogin)
	list_commands.register("register", handleRegister)

	args := os.Args
	if len(args) < 2 {
		log.Fatal("Cannot execute any command ! There is less than two arguments")
	}
	cmd := command{args[1], args[2:]}
	err = list_commands.run(&app, cmd)
	if err != nil {
		log.Fatal(err)
	}
}
