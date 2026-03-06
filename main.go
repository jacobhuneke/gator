package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/jacobhuneke/gator/internal/config"
	"github.com/jacobhuneke/gator/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	//opens db, sets queries
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dbQueries := database.New(db)

	//gets data from config file
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//inits state struct which stores current cfg and db
	s := state{
		cfg: cfg,
		db:  dbQueries,
	}

	//makes a map of commands
	cmds := commands{
		funcMap: make(map[string]func(*state, command) error),
		descMap: make(map[string]string),
	}

	s.cmds = &cmds

	e := registerCommands(*s.cmds)
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}

	//gets passed arguments
	args := os.Args

	//creates command
	newCmd := command{
		name: args[1],
		args: args[2:],
	}

	//runs command
	runErr := cmds.run(&s, newCmd)
	if runErr != nil {
		fmt.Println(runErr)
		os.Exit(1)
	}
}
