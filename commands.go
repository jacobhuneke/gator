package main

import (
	"errors"

	"github.com/jacobhuneke/gator/internal/config"
	"github.com/jacobhuneke/gator/internal/database"
)

type state struct {
	cfg  *config.Config
	db   *database.Queries
	cmds *commands
}

type command struct {
	name        string
	args        []string
	description string
}

type commands struct {
	funcMap map[string]func(*state, command) error
	descMap map[string]string
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

func (c *commands) register(name, description string, f func(*state, command) error) error {
	_, ok := c.funcMap[name]
	_, ok2 := c.descMap[name]
	if ok && ok2 {
		return errors.New("command already registered")
	} else {
		c.funcMap[name] = f
		c.descMap[name] = description
		return nil
	}
}

func registerCommands(cmds commands) error {
	//registers login, register command
	e := cmds.register("login", "Sets the current user to the given input. Takes a name input.", handlerLogin)
	if e != nil {
		return e
	}

	e = cmds.register("register", "Registers a user by a given input and sets them as the current user. Takes a name input.", handlerRegister)
	if e != nil {
		return e
	}

	e = cmds.register("reset", "Resets the users and feeds DBs. No args required.", handlerReset)
	if e != nil {
		return e
	}

	e = cmds.register("users", "Lists all users in the DB and identifies the current one. No args required.", handlerUsers)
	if e != nil {
		return e
	}

	e = cmds.register("agg", "Collects feed. Takes a time interval requirement arg.", handlerAgg)
	if e != nil {
		return e
	}

	e = cmds.register("addfeed", "Adds a new feed to the current user. Takes a name and url args.", middlewareLoggedIn(handlerAddfeed))
	if e != nil {
		return e
	}

	e = cmds.register("feeds", "Prints all feeds to the console. Takes no args.", handlerFeeds)
	if e != nil {
		return e
	}

	e = cmds.register("follow", "Creates a new feed follow record, which stores information for all the feeds the current user follows. Takes a url arg.", middlewareLoggedIn(handlerFollow))
	if e != nil {
		return e
	}

	e = cmds.register("following", "Prints all the names of the feeds the current user is following. Takes no args.", middlewareLoggedIn(handlerFollowing))
	if e != nil {
		return e
	}

	e = cmds.register("unfollow", "Unfollows a feed for the current user and removes it from the feed following table. Takes a url arg.", middlewareLoggedIn(handlerUnfollow))
	if e != nil {
		return e
	}

	e = cmds.register("browse", "Prints posts information to the terminal. Takes an optional limit arg. Default is 2.", middlewareLoggedIn(handlerBrowse))
	if e != nil {
		return e
	}
	e = cmds.register("help", "Lists descriptions of all commands. Takes no args.", handlerHelp)
	if e != nil {
		return e
	}
	return nil
}
