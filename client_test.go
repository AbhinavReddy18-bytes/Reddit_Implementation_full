package client_test

import (
	"project4/client"
	"project4/engine"
	"testing"
)

func TestClientRegistration(t *testing.T) {
	e := engine.NewEngine()
	user := client.NewClient("user1", e)

	err := user.Register()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestClientJoinSubreddit(t *testing.T) {
	e := engine.NewEngine()
	user := client.NewClient("user1", e)

	user.Register()
	e.CreateSubreddit("test_subreddit")

	err := user.JoinSubreddit("test_subreddit")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestClientPosting(t *testing.T) {
	e := engine.NewEngine()
	user := client.NewClient("user1", e)

	user.Register()
	e.CreateSubreddit("test_subreddit")

	postID, err := user.PostInSubreddit("test_subreddit", "This is a test post.")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if postID != 1 {
		t.Errorf("expected post ID to be 1, got %d", postID)
	}
}

func TestClientMessaging(t *testing.T) {
	e := engine.NewEngine()
	user1 := client.NewClient("user1", e)
	user2 := client.NewClient("user2", e)

	user1.Register()
	user2.Register()

	err := user1.SendMessage("user2", "Hello from user1!")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	user2.ListMessages()
}
