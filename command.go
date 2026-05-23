package main

import (
	"context"
	"errors"
	"fmt"
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

func CreateCommands() commands {
	ret := commands{handlers: make(map[string]func(*state, command) error)}

	ret.register("login", handlerLogin)
	ret.register("register", handlerRegister)
	ret.register("reset", handlerReset)
	ret.register("users", handlerUsers)

	return ret
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

	user, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("User: %v, not found. err: %v", cmd.args[0], err)
	}

	err = s.config.SetUser(user.Name)
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

func handlerReset(s *state, cmd command) error {
	return s.db.ResetUsers(context.Background())
}

func handlerUsers(s *state, cmd command) error {
	names, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, name := range names {
		if s.config.CurrentUserName == name {
			fmt.Println("*", name, "(current)")
		} else {
			fmt.Println("*", name)
		}
	}

	return nil
}
