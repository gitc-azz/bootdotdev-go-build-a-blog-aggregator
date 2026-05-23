package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/gitc-azz/bootdotdev-go-build-a-blog-aggregator/internal/database"
	"github.com/google/uuid"
)

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (self *commands) run(s *state, cmd command) error {
	err := self.handlers[cmd.name](s, cmd)

	return err
}

func (self *commands) register(name string, handler func(*state, command) error) {
	self.handlers[name] = handler
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("Login command expect one argument, the username")
	}

	err := s.config.SetUser(cmd.args[0])
	if err != nil {
		return err
	}

	log.Println("The user has been set to", s.config.CurrentUserName)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("Register command expect one argument, the username")
	}

	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}

	user_db, err := s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		return err
	}

	err = s.config.SetUser(user_db.Name)
	if err != nil {
		return err
	}

	log.Println("The user", user_db.Name, "was created")

	return nil
}
