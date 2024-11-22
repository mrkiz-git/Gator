package main

import (
	"fmt"
	"log"
	"mrkiz-git/gator/internal/config"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	cfg.SetUser("mark")
	cfg, err = config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("%+v\n", *cfg)

}
