package main

import (
	"log"
	"os"

	"github.com/gitc-azz/bootdotdev-go-build-a-blog-aggregator/internal/config"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatal("Usage: gator <command> <argument>")
	}

	conf, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	st := state{config: &conf}

	cmds := commands{handlers: make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)

	loginCmd := command{name: os.Args[1], args: []string{os.Args[2]}}

	err = cmds.run(&st, loginCmd)
	if err != nil {
		log.Fatal(err)
	}

	// confUpdated, err := config.Read()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(confUpdated)
}

type state struct {
	config *config.Config
}
