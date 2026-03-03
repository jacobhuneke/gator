package main

import (
	"errors"
	"fmt"

	"github.com/jacobhuneke/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("the login handler expects a username argument")
	}

	err := s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	} else {
		fmt.Println("The user has been set")
	}

	return nil
}

type commands struct {
	funcMap map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	exe, ok := c.funcMap[cmd.name]
	if ok {
		exe(s, cmd)
		return nil
	} else {
		return errors.New("command is not valid")
	}
}

func (c *commands) register(name string, f func(*state, command) error) error {
	_, ok := c.funcMap[name]
	if ok {
		return errors.New("command already registered")
	} else {
		c.funcMap[name] = f
		return nil
	}
}
