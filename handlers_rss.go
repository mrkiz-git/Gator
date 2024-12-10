package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"mrkiz-git/gator/internal/database"
	"mrkiz-git/gator/internal/rss"
	"time"

	"github.com/google/uuid"
)

const rssFeedURL = "https://www.wagslane.dev/index.xml"

func handlFetchRSSFeed(s *state, cmd command) error {

	if len(cmd.Args) != 1 {
		return fmt.Errorf("wrong number of parameters")
	}

	time_between_reqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}
	tiker := time.NewTicker(time_between_reqs)

	log.Printf("Collecting feeds every %v", time_between_reqs)

	for ; ; <-tiker.C {
		nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
		if err != nil {
			return err
		}
		log.Printf("*****  Items form RSS FEED: %s saved to DB", nextFeed.Name)
		rssFeed, err := rss.FetchFeed(context.Background(), nextFeed.Url)
		if err != nil {
			return err
		}

		for _, item := range rssFeed.Channel.Item {

			pubTime, err := time.Parse(time.RFC1123, item.PubDate)
			if err != nil {
				return err
			}
			postData := database.CreatePostParams{
				ID:          uuid.New(),
				FeedID:      nextFeed.ID,
				Title:       item.Title,
				Description: item.Description,
				Url:         item.Link,
				PublishedAt: pubTime,
			}
			insertResponse, err := s.db.CreatePost(context.Background(), postData)

			if err != nil {
				return err
			}
			log.Printf("Post %s inserted the data base", insertResponse.Title)
		}
		update_time := sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		prms := database.MarkFeedFetchedParams{LastFetchedAt: update_time, ID: nextFeed.ID}

		_, err = s.db.MarkFeedFetched(context.Background(), prms)

		if err != nil {
			return err
		}

	}
	return nil
}

func haddleAddRSSFeed(s *state, cmd command, user database.User) error {

	if len(cmd.Args) != 2 {
		return fmt.Errorf("Wrong number of paramters")
	}

	userUUID := uuid.NullUUID{
		UUID:  user.ID,
		Valid: true, // Indicate that this UUID is valid (not NULL)
	}

	feedName := cmd.Args[0]
	feedURL := cmd.Args[1]
	newFeed := database.CreateFeedParams{
		ID:        uuid.New(),
		UserID:    userUUID,
		Name:      feedName,
		Url:       feedURL,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// check if feed exists
	feeds, err := s.db.GetFeedByURL(context.Background(), feedURL)

	if err != nil {
		log.Printf("error geting feed by id")
		return err
	}
	if len(feeds) > 1 || len(feeds) < 0 {
		return fmt.Errorf("unxepected amount of feeds with the same url")
	}
	if len(feeds) == 0 {

		feed, err := s.db.CreateFeed(context.Background(), newFeed)
		if err != nil {
			return err
		}
		printFeed(feed)
	}

	handleFollowFeed(s, command{Name: "follow", Args: cmd.Args[len(cmd.Args)-1:]}, user)

	return nil

}
func handleListFeeds(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("zero argumetns expected")
	}

	feeds, err := s.db.ListFeeds(context.Background())

	if err != nil {
		log.Printf("error loading feeds list")
		return err
	}

	for _, feed := range feeds {
		fmt.Println("Feed Details:")
		fmt.Printf("**	ID: %s\n", feed.ID)
		fmt.Printf("**	User Name: %s\n", feed.UserName)
		fmt.Printf("**	User ID: %s\n", feed.UserID.UUID)
		fmt.Printf("**	Name: %s\n", feed.Name)
		fmt.Printf("**	URL: %s\n", feed.Url)
		fmt.Printf("**	Created At: %s\n", feed.CreatedAt.Format(time.RFC3339))
		fmt.Printf("**	Updated At: %s\n", feed.UpdatedAt.Format(time.RFC3339))

	}

	return nil
}

func handleFollowFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("one argument expected")
	}

	feeds, err := s.db.GetFeedByURL(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}

	if len(feeds) != 1 {
		return fmt.Errorf("unexpected  amount of feeds recived")
	}

	feedUUID := feeds[0].ID

	newFeedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    feedUUID,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}

	newfeedfollow, err := s.db.CreateFeedFollow(context.Background(), newFeedFollow)
	if err != nil {
		return err
	}

	fmt.Printf("User: %s started follow RSS feed: %s", newfeedfollow.UserName, newfeedfollow.FeedName)
	return nil
}

func handleGetFollowedFeeds(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("wrong number of parameters")
	}

	followedFeeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	for _, feed := range followedFeeds {
		fmt.Printf("** %s\n", feed.FeedName)
	}
	return nil
}

func handleUnfollowFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Wrong Number of Arguments")
	}
	feedUrl := cmd.Args[0]
	feeds, err := s.db.GetFeedByURL(context.Background(), feedUrl)
	if err != nil {
		return err
	}

	unfolowParams := database.UnfolowFeedParams{UserID: user.ID,
		FeedID: feeds[0].ID,
	}

	feed, err := s.db.UnfolowFeed(context.Background(), unfolowParams)

	if err != nil {
		return err
	}
	log.Printf("Feed %v was removed for user %v", feed.FeedID, feed.UserID)
	return nil
}

func printFeed(feed database.CreateFeedRow) {
	fmt.Println("Feed Details:")
	fmt.Printf("**	ID: %s\n", feed.ID)

	if feed.UserID.Valid {
		fmt.Printf("** User ID: %s\n", feed.UserID.UUID)
	} else {
		fmt.Println("**	User ID: NULL")
	}

	fmt.Printf("**	Name: %s\n", feed.Name)
	fmt.Printf("**	URL: %s\n", feed.Url)
	fmt.Printf("**	Created At: %s\n", feed.CreatedAt.Format(time.RFC3339))
	fmt.Printf("**	Updated At: %s\n", feed.UpdatedAt.Format(time.RFC3339))
}

func scrapeFeeds(s *state) {

	feedToScrape, _ := s.db.GetNextFeedToFetch(context.Background())

	update_time := sql.NullTime{Time: time.Now(), Valid: true}
	markParams := database.MarkFeedFetchedParams{LastFetchedAt: update_time, ID: feedToScrape.ID}

	_, err := s.db.MarkFeedFetched(context.Background(), markParams)

	if err != nil {
		log.Fatal("failed to updaed feed date")
	}

	feed, err := rss.FetchFeed(context.Background(), feedToScrape.Url)
	if err != nil {
		log.Fatal("failed to scrape feed")
	}

	fmt.Println(feed)

}
