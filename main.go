package main

import (
	"log"
	"os"

	_ "github.com/lib/pq" // indirectly required to use postgres

	"github.com/gitc-azz/bootdotdev-go-build-a-blog-aggregator/repl"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: gator <command> [argument]")
	}

	state, cmds := repl.Init()

	err := cmds.Run(state, repl.NewCommand(os.Args[1], os.Args[2:]))
	if err != nil {
		log.Fatal(err)
	}
}
