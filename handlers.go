package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jacobhuneke/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			fmt.Println("no users to call command on")
			return err
		}
		return handler(s, cmd, user)
	}
}

func handlerRegister(s *state, cmd command) error {
	//exits early if user exists, or no name is passed
	if len(cmd.args) == 0 {
		return fmt.Errorf("the register handler expects a name argument")
	}

	if _, err := s.db.GetUser(context.Background(), cmd.args[0]); err == nil {
		fmt.Printf("the user %s already exists\n", cmd.args[0])
		os.Exit(1)
	}

	//create user parameters
	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}

	user, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return err
	}

	//sets the current user, updates successfully
	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("The user was successfully created")
	fmt.Println(user)
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("the login handler expects a username argument")
	}
	if _, err := s.db.GetUser(context.Background(), cmd.args[0]); err != nil {
		fmt.Printf("the user %s does not exist\n", cmd.args[0])
		fmt.Println(err)
		os.Exit(1)
	}

	err := s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return err
	} else {
		fmt.Println("The user has been set")
	}

	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("the deletion was successful")
	os.Exit(0)
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		if user.Name != s.cfg.CurrentUserName {
			fmt.Printf(" * %v\n", user.Name)
		} else {
			fmt.Printf(" * %v (current)\n", user.Name)
		}
	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("must provide time between reqs")
	}

	dur, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Collecting feeds every %v\n", dur)
	ticker := time.NewTicker(dur)
	scrapeFeeds(s)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func handlerAddfeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("not enough arguments provided")
	}

	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return err
	}

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feedParams.ID,
	}
	_, err = s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return err
	}

	fmt.Printf("ID: %v\n", feed.ID)
	fmt.Printf("CreatedAt: %v\n", feed.CreatedAt)
	fmt.Printf("UpdatedAt: %v\n", feed.UpdatedAt)
	fmt.Printf("Name: %v\n", feed.Name)
	fmt.Printf("Url: %v\n", feed.Url)
	fmt.Printf("UserID: %v\n", feed.UserID)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Println(feed.Name)
		fmt.Println(feed.Url)
		username, err := s.db.GetUserName(context.Background(), feed.UserID)
		if err != nil {
			return err
		}
		fmt.Println(username)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("must provide a url")
	}

	feed, err := s.db.GetFeedFromURL(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feed_follows, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Println(feed_follows.FeedName)
	fmt.Println(feed_follows.UserName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {

	feedFollowing, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	for _, feedFollow := range feedFollowing {
		feed, err := s.db.GetFeedFromID(context.Background(), feedFollow.FeedID)
		if err != nil {
			return err
		}
		fmt.Println(feed.Name)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("must provide a url")
	}

	feed, err := s.db.GetFeedFromURL(context.Background(), cmd.args[0])
	if err != nil {
		return err
	}

	params := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	e := s.db.DeleteFeedFollow(context.Background(), params)
	if e != nil {
		return e
	}

	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int32
	if len(cmd.args) == 1 {
		l64, err := strconv.ParseInt(cmd.args[0], 10, 32)
		if err != nil {
			return err
		}
		limit = int32(l64)
	} else {
		limit = 2
	}
	params := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	}
	rows, err := s.db.GetPostsForUser(context.Background(), params)
	if err != nil {
		return err
	}

	for _, row := range rows {
		fmt.Printf("--- %s ---\n", row.Title)
		fmt.Printf("Source: %s\n", row.Url)
		if row.Description != "" {
			fmt.Printf("Description: %s\n", row.Description)
		}
		fmt.Println("-------------------------------------")
	}
	return nil
}

func handlerHelp(s *state, cmd command) error {
	for name, desc := range s.cmds.descMap {
		fmt.Printf("%s: %s\n", name, desc)
	}
	return nil
}
