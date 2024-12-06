package main

import (
	"context"
	"fmt"
	"time"

	"mrkiz-git/gator/internal/database"

	"github.com/google/uuid"
)

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <name>", cmd.Name)
	}

	name := cmd.Args[0]

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
	})
	if err != nil {
		return fmt.Errorf("couldn't create user: %w", err)
	}

	err = s.cfg.SetUser(user.Name, user.ID)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User created successfully:")
	printUser(user)
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	name := cmd.Args[0]

	user, err := s.db.GetUserByName(context.Background(), name)
	if err != nil {
		return fmt.Errorf("couldn't find user: %w", err)
	}

	err = s.cfg.SetUser(name, user.ID)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User switched successfully!")
	return nil
}

func handlerRestDB(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("zero arguments expected")
	}
	_, err := s.db.ResetDataBase(context.Background())
	if err != nil {
		return fmt.Errorf("failed to reset database, %w", err)
	}

	return nil
}

func handlerGetUsers(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("zero arguments expected")
	}

	res, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to fetch user list %w", err)
	}
	for _, r := range res {
		if r.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", r.Name)
		} else {
			fmt.Printf("* %s\n", r.Name)
		}

	}
	return nil
}

func printUser(user database.User) {
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)
}
