package main

import (
	"errors"
	"log"
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
