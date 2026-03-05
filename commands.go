package main

import (
	"errors"

	"github.com/jacobhuneke/gator/internal/config"
	"github.com/jacobhuneke/gator/internal/database"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

type command struct {
	name string
	args []string
}

type commands struct {
	funcMap map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	exe, ok := c.funcMap[cmd.name]
	if ok {
		err := exe(s, cmd)
		if err != nil {
			return err
		}
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

func registerCommands(cmds commands) error {
	//registers login, register command
	e := cmds.register("login", handlerLogin)
	if e != nil {
		return e
	}

	e = cmds.register("register", handlerRegister)
	if e != nil {
		return e
	}

	e = cmds.register("reset", handlerReset)
	if e != nil {
		return e
	}

	e = cmds.register("users", handlerUsers)
	if e != nil {
		return e
	}

	e = cmds.register("agg", handlerAgg)
	if e != nil {
		return e
	}

	e = cmds.register("addfeed", middlewareLoggedIn(handlerAddfeed))
	if e != nil {
		return e
	}

	e = cmds.register("feeds", handlerFeeds)
	if e != nil {
		return e
	}

	e = cmds.register("follow", middlewareLoggedIn(handlerFollow))
	if e != nil {
		return e
	}

	e = cmds.register("following", middlewareLoggedIn(handlerFollowing))
	if e != nil {
		return e
	}
	return nil
}
