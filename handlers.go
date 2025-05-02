package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gator/internal/config"
	"gator/internal/database"
	"time"
)

type state struct {
	config *config.Config
	db     *database.Queries
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.parameters) != 1 {
		return errors.New("login command expects 1 parameter: [username]")
	}

	username := cmd.parameters[0]

	user, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return errors.New("username doesn't exist")
	}

	s.config.Current_user_name = user.Name
	err = s.config.SetUser()
	if err != nil {
		return err
	}

	fmt.Printf("Login successfully as [%s]\n", user.Name)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.parameters) != 1 {
		return errors.New("register command expects 1 parameters: [username]")
	}

	username := cmd.parameters[0]

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      username,
	})

	if err != nil {
		return err
	}

	fmt.Printf(`User created successfully
	UUID: %s
	CreatedAt: %s
	UpdatedAt: %s
	Name: %s\n
`, user.ID, user.CreatedAt.Local().String(), user.UpdatedAt.Local().String(), user.Name)

	s.config.Current_user_name = user.Name
	s.config.SetUser()

	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.parameters) != 0 {
		return errors.New("reset command expects 0 parameters")
	}

	err := s.db.DatabaseReset(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("Reset successfully.")

	return nil
}

func handlerUsers(s *state, cmd command) error {
	if len(cmd.parameters) != 0 {
		return errors.New("users command expects 0 parameters")
	}

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	currentUser := s.config.Current_user_name

	for _, user := range users {
		fmt.Printf("* %s", user.Name)
		if user.Name == currentUser {
			fmt.Printf(" (current)")
		}
		fmt.Println()
	}

	return nil
}

func handlerAggregate(s *state, cmd command) error {
	if len(cmd.parameters) != 0 {
		return errors.New("agg command expects 0 parameters")
	}

	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Print(feed, "\n")

	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.parameters) != 2 {
		return errors.New("addfeed command expects 2 parameters: [name] [url]")
	}

	user, err := s.db.GetUser(context.Background(), s.config.Current_user_name)
	if err != nil {
		return err
	}
	feedName := cmd.parameters[0]
	url := cmd.parameters[1]

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		Name:      sql.NullString{String: feedName, Valid: true},
		Url:       sql.NullString{String: url, Valid: true},
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	fmt.Printf(`Feed created successfully
	Name: %s
	Url: %s
	CreatedAt: %s
	UpdatedAt: %s
	User: %s
`, feed.Name.String, feed.Url.String, feed.CreatedAt.Local().String(), feed.UpdatedAt.Local().String(), user.Name)

	return nil
}

func handlerShowFeeds(s *state, cmd command) error {
	if len(cmd.parameters) != 0 {
		return errors.New("feeds command expects 0 parameters")
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return err
	}

	for _, feed := range feeds {
		owner, err := s.db.GetUserFromId(context.Background(), feed.UserID)
		if err != nil {
			return err
		}

		fmt.Printf("Feed name: %s\n", feed.Name.String)
		fmt.Printf("Url: %s\n", feed.Url.String)
		fmt.Printf("Created by: %s\n\n", owner.Name)
	}

	return nil
}
