package main

import (
	"fmt"

	"github.com/jacobhuneke/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	e := cfg.SetUser("Jake")
	if e != nil {
		fmt.Println(e)
	}

	c, e2 := config.Read()
	if e2 != nil {
		fmt.Println(e2)
	}
	fmt.Println(c)
}
