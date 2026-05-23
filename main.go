package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/gitc-azz/bootdotdev-go-build-a-blog-aggregator/internal/config"
	"github.com/gitc-azz/bootdotdev-go-build-a-blog-aggregator/internal/database"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: gator <command> [argument]")
	}

	conf, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	st := state{config: &conf}

	db, err := sql.Open("postgres", st.config.DbUrl)
	if err != nil {
		log.Fatal(err)
	}

	st.db = database.New(db)

	cmds := CreateCommands()

	cmd := command{name: os.Args[1], args: os.Args[2:]}

	err = cmds.run(&st, cmd)
	if err != nil {
		log.Fatal(err)
	}
}

type state struct {
	config *config.Config
	db     *database.Queries
}
