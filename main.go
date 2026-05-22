package main

import (
	"fmt"
	"log"

	"github.com/gitc-azz/bootdotdev-go-build-a-blog-aggregator/internal/config"
)

func main() {
	conf, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	conf.SetUser("azz")

	confUpdated, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(confUpdated)
}

type state struct {
	config *config.Config
}
