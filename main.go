package main

import (
	"fmt"
	"log"
	"mrkiz-git/gator/internal/config"
	"os"
)

type state struct {
	cfg *config.Config
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	// cretae a pointer to state, every time programState w be passed, it will be passed to *state
	programState := &state{
		cfg: cfg,
	}

	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)

	if len(os.Args) < 2 {
		fmt.Println("Usage: cli <command> [args...]")
		return
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.run(programState, command{Name: cmdName, Args: cmdArgs})
	if err != nil {
		log.Fatal(err)
	}
}
