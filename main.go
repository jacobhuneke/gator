package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/jacobhuneke/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	s := state{
		cfg: cfg,
	}

	cmds := commands{
		funcMap: make(map[string]func(*state, command) error),
	}

	e := cmds.register("login", handlerLogin)
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}

	args := os.Args
	if len(args) <= 2 {
		err := errors.New("not enough arguments provided")
		fmt.Println(err)
		os.Exit(1)
	}

	newCmd := command{
		name: args[1],
		args: args[2:],
	}

	err = cmds.run(&s, newCmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
