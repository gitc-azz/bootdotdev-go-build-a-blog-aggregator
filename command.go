package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gitc-azz/bootdotdev-go-build-a-blog-aggregator/internal/database"
	"github.com/google/uuid"
)

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func CreateCommands() commands {
	ret := commands{handlers: make(map[string]func(*state, command) error)}

	ret.register("login", handlerLogin)
	ret.register("register", handlerRegister)
	ret.register("reset", handlerReset)
	ret.register("users", handlerUsers)
	ret.register("agg", handlerAgg)
	ret.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	ret.register("feeds", handlerFeeds)
	ret.register("follow", middlewareLoggedIn(handlerFollow))
	ret.register("following", middlewareLoggedIn(handlerFollowing))
	ret.register("unfollow", middlewareLoggedIn(handlerUnfollow))

	return ret
}

func (self *commands) run(s *state, cmd command) error {
	if _, ok := self.handlers[cmd.name]; !ok {
		return fmt.Errorf("%v is a unknown command", cmd.name)
	}
	err := self.handlers[cmd.name](s, cmd)

	return err
}

func (self *commands) register(name string, handler func(*state, command) error) {
	self.handlers[name] = handler
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("Login command expect one argument, the username")
	}

	user, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("User: %v, not found. err: %v", cmd.args[0], err)
	}

	err = s.config.SetUser(user.Name)
	if err != nil {
		return err
	}

	log.Println("The user has been set to", s.config.CurrentUserName)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return errors.New("Register command expect one argument, the username")
	}

	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}

	user_db, err := s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		return err
	}

	err = s.config.SetUser(user_db.Name)
	if err != nil {
		return err
	}

	log.Println("The user", user_db.Name, "was created")

	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to reset users table, err= %v", err)
	}

	err = s.db.ResetFeedFollows(context.Background())
	if err != nil {
		return fmt.Errorf("failed to reset feed_follows table, err= %v", err)
	}

	return s.db.ResetFeed(context.Background())
}

func handlerUsers(s *state, cmd command) error {
	names, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, name := range names {
		if s.config.CurrentUserName == name {
			fmt.Println("*", name, "(current)")
		} else {
			fmt.Println("*", name)
		}
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return errors.New("agg command takes no arguments")
	}

	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Println(feed)

	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("%v command expects two arguments: name and url of the feed", cmd.name)
	}
	feedName, feedURL := cmd.args[0], cmd.args[1]

	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedURL,
	}
	feed, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return fmt.Errorf("Failed to add feed to db, err: %v", err)
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("Failed to create the corresponding feed_follow, err: %v", err)
	}

	fmt.Println("feed:", feedFollow.FeedName, ", added and followed by", feedFollow.UserName)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("%v command takes no arguments", cmd.name)
	}

	feedsWithUserName, err := s.db.GetFeedsWithUserName(context.Background())
	if err != nil {
		return err
	}

	for _, fwun := range feedsWithUserName {
		fmt.Println("feed:", fwun.Name, ", url:", fwun.Url, "user:", fwun.Username)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("%v command takes exactly one argument, an url to RSS feed", cmd.name)
	}
	urlFeed := cmd.args[0]

	feed, err := s.db.GetFeed(context.Background(), urlFeed)
	if err != nil {
		return fmt.Errorf("feed at %v not found, err= %v", urlFeed, err)
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create a feed_follow, err= %v", err)
	}

	fmt.Println("feed:", feedFollow.FeedName, "is now followed by", feedFollow.UserName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("%v command expects zero argument", cmd.name)
	}

	fffu, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return fmt.Errorf("failed to get the RSS feed followed by %v, with err= %v", user.Name, err)
	}

	if len(fffu) == 0 {
		fmt.Println(user.Name, "don't follow any RSS feed !")
		return nil
	}

	fmt.Println(user.Name, "is subscribed to the following feed:")

	for _, ff := range fffu {
		fmt.Println(" *", ff.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("%v command expect one argument: url to a RSS feed", cmd.name)
	}
	urlRSS := cmd.args[0]

	feed, err := s.db.GetFeed(context.Background(), urlRSS)
	if err != nil {
		return fmt.Errorf("failed to get RSS feed: %v", urlRSS)
	}

	err = s.db.RmFeedFollow(context.Background(), database.RmFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to remove the RSS feed: %v, for the user: %v", feed.Name, user.Name)
	}

	fmt.Println(user.Name, "has unfollowed:", feed.Url)

	return nil
}

func middlewareLoggedIn(handler func(*state, command, database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
		if err != nil {
			return fmt.Errorf("Failed to fetch %v user from db, are you registered?, err: %v", s.config.CurrentUserName, err)
		}

		return handler(s, cmd, user)
	}
}
