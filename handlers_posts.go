package main

import (
	"context"
	"fmt"
	"strconv"
)

func handlBrowsPosts(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Wrong number of parameters")
	}

	limit, err := strconv.Atoi(cmd.Args[0])
	if err != nil {
		return err
	}

	queryResult, err := s.db.GetAllPosts(context.Background(), int32(limit))
	if err != nil {
		return err
	}

	for _, val := range queryResult {
		fmt.Printf("Title: %s\nURL:%s\n ", val.Title, val.Url)
	}
	return nil

}
