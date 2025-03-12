package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/willthefoollearn/gator/internal/database"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return errors.New("please submit with a single username")
	}

	user, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		fmt.Println("User wasn't found")
		os.Exit(1)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("Username has been set to %s\n", cmd.args[0])

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.name) == 0 {
		return errors.New("please submit with command name")
	}

	if len(cmd.args) == 0 {
		return errors.New("please submit with a username")
	}

	var arg = database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	}

	newUser, err := s.db.CreateUser(context.Background(), arg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = s.cfg.SetUser(newUser.Name)
	if err != nil {
		return err
	}

	fmt.Printf("User %s was created with id: %v \n", newUser.Name, newUser.ID)

	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.ResetUser(context.Background())
	if err != nil {
		fmt.Println("Unable to reset database")
		os.Exit(1)
	}

	fmt.Println("Database was reset")

	return nil
}

func handlerList(s *state, cmd command) error {
	users, err := s.db.ListUsers(context.Background())
	if err != nil {
		fmt.Println("Unable to find users")
		os.Exit(1)
	}

	for _, user := range users {
		name := user.Name

		if s.cfg.CurrentUserName == name {
			name += " (current)"
		}

		fmt.Printf("* %s\n", name)
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("usage: %s <time>", cmd.name)
	}

	dur, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("couldn't determine duration %s", err)
	}

	fmt.Printf("Collecting feeds every %v\n", dur)

	ticker := time.NewTicker(dur)

	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Printf("Unable to pull feed: %v", err)
		return
	}

	scrapeFeed(s.db, feed)
}

func scrapeFeed(db *database.Queries, feeds database.Feed) {
	feed, err := db.MarkFeedFetched(context.Background(), feeds.ID)
	if err != nil {
		log.Printf("Unable to mark feed: %v", err)
		return
	}

	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("Unable to pull feed data: %v", err)
		return
	}

	for _, data := range feedData.Channel.Item {
		convertedTime, err := time.Parse(time.RFC1123Z, data.PubDate)
		if err != nil {
			log.Printf("Unable to convert time: %v", err)
		}

		var arg = database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       data.Title,
			Url:         data.Link,
			Description: data.Description,
			PublishedAt: convertedTime,
			FeedID:      feed.ID,
		}

		createdPost, err := db.CreatePost(context.Background(), arg)
		if err != nil {
			log.Printf("Post wasn't created: %v", err)
		}

		fmt.Printf("New post added to database: %s\n", createdPost.Title)
	}

	log.Printf("A feed was found for %s with %d articles\n\n", feed.Name, len(feedData.Channel.Item))
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) < 2 {
		fmt.Println("Not enough arguments, you dingus!")
		os.Exit(1)
	}

	name := cmd.args[0]
	url := cmd.args[1]

	var feed = database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	}

	newFeed, err := s.db.CreateFeed(context.Background(), feed)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cmd.args[0] = cmd.args[1]

	err = handlerFollow(s, cmd, user)
	if err != nil {
		return errors.New("couldn't add to feed follow")
	}

	fmt.Printf("%+v\n", newFeed)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.ListFeeds(context.Background())
	if err != nil {
		fmt.Println("Can't get all the feeds")
		return err
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found")
		return nil
	}

	for _, feed := range feeds {
		userName, err := s.db.FeedsUser(context.Background(), feed.UserID)
		if err != nil {
			fmt.Println("Username can't be found")
			return err
		}

		fmt.Printf("Feed Name: %s\n", feed.Name)
		fmt.Printf("Feed URL: %s\n", feed.Url)
		fmt.Printf("Feed's User: %s\n\n", userName.Name)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	feed, err := s.db.FeedFromUrl(context.Background(), cmd.args[0])
	if err != nil {
		return errors.New("current feed not found from URL")
	}

	var feedFollow = database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	createdFeed, err := s.db.CreateFeedFollows(context.Background(), feedFollow)
	if err != nil {
		return errors.New("feed could not be created")
	}

	fmt.Printf("Current user: %s\n", createdFeed.UserName)
	fmt.Printf("Current user: %s\n", createdFeed.FeedName)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("usage: %s", cmd.name)
	}

	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return errors.New("current feed couldn't be found")
	}

	for _, feed := range feeds {
		fmt.Printf("%s\n", feed.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) > 1 {
		return fmt.Errorf("usage: %s <url>", cmd.name)
	}

	feed, err := s.db.FeedFromUrl(context.Background(), cmd.args[0])
	if err != nil {
		return errors.New("feed wasn't found")
	}

	var args = database.UnfollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	err = s.db.Unfollow(context.Background(), args)
	if err != nil {
		return errors.New("couldn't unfollow")
	}

	log.Printf("Successfully unfollowed %s", feed.Name)

	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	limitCount := 2
	var err error

	if len(cmd.args) == 1 {
		limitCount, err = strconv.Atoi(cmd.args[0])
		if err != nil {
			return errors.New("unable to get limit")
		}
	}

	var arg = database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limitCount),
	}

	posts, err := s.db.GetPostsForUser(context.Background(), arg)
	if err != nil {
		return errors.New("couldn't grab posts")
	}

	for _, post := range posts {
		fmt.Printf("%s: \n%s\n", post.Title, post.Description)
	}

	return nil
}
