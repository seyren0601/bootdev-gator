package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/seyren0601/bootdev-gator/internal/config"
	"github.com/seyren0601/bootdev-gator/internal/database"
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
	fmt.Printf("Logged in as [%s]\n", user.Name)

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
	if len(cmd.parameters) != 1 {
		return errors.New("agg command expects 1 parameters: [time_between_reqs]")
	}

	time_between_reqs, err := time.ParseDuration(cmd.parameters[0])
	if err != nil {
		return errors.New("can't parse parameter into duration")
	}

	ticker := time.NewTicker(time_between_reqs)
	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			return err
		}
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.parameters) != 2 {
		return errors.New("addfeed command expects 2 parameters: [name] [url]")
	}

	feedName := cmd.parameters[0]
	url := cmd.parameters[1]

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		Name:      sql.NullString{String: feedName, Valid: true},
		Url:       url,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		UserID:    user.ID,
		FeedID:    feed.ID,
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
[Automatically] You (%s) have followed this feed.
`, feed.Name.String, feed.Url, feed.CreatedAt.Local().String(), feed.UpdatedAt.Local().String(), user.Name, user.Name)

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
		fmt.Printf("Url: %s\n", feed.Url)
		fmt.Printf("Created by: %s\n\n", owner.Name)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.parameters) != 1 {
		return errors.New("follow command expects 1 parameters: [url]")
	}

	url := cmd.parameters[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return err
	}

	follow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		FeedID:    feed.ID,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	fmt.Println("Follow successfully.")
	fmt.Printf("Feed: %s\n", follow.Feedname.String)
	fmt.Printf("User: %s\n", follow.Username)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.parameters) != 0 {
		return errors.New("following command expects 0 parameters")
	}

	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	fmt.Println("Current feed follows:")
	for _, feed := range follows {
		fmt.Printf("\t* %s\n", feed.FeedName.String)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.parameters) != 1 {
		return errors.New("unfollow command expects 1 parameters: [url]")
	}

	url := cmd.parameters[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return err
	}

	err = s.db.DeleteFeedFollowForUser(context.Background(), database.DeleteFeedFollowForUserParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return err
	}

	fmt.Printf("You (%s) have unfollowed '%s' at '%s'\n", user.Name, feed.Name.String, feed.Url)

	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	var limit int = 2
	var err error
	if len(cmd.parameters) > 0 {
		limit, err = strconv.Atoi(cmd.parameters[0])
		if err != nil {
			return err
		}
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return err
	}

	if len(posts) < 1 {
		return errors.New("you aren't following any feed")
	}

	// print from posts[limit - 1] to posts[0]
	// so that the last post printed (on console) is the most recent
	for i := limit - 1; i >= 0; i-- {
		post := posts[i]
		fmt.Printf(`Source: %s
	Title: %s
	Published at: %s

`, post.Source.String, post.Title, post.PublishedAt.Local().Format(time.ANSIC))
	}

	return nil
}
