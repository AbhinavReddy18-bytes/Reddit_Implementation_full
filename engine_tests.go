package engine_tests

import (
	"project4/engine"
	"testing"
)

func TestRegisterUser(t *testing.T) {
	e := engine.NewEngine()

	err := e.RegisterUser("test_user")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = e.RegisterUser("test_user")
	if err == nil {
		t.Errorf("expected error for duplicate user, got none")
	}
}

func TestSubredditManagement(t *testing.T) {
	e := engine.NewEngine()

	err := e.CreateSubreddit("test_subreddit")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = e.RegisterUser("test_user")
	if err != nil {
		t.Fatalf("unexpected error during user registration: %v", err)
	}
	err = e.JoinSubreddit("test_user", "test_subreddit")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = e.LeaveSubreddit("test_user", "test_subreddit")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestPostInSubreddit(t *testing.T) {
	e := engine.NewEngine()

	e.CreateSubreddit("test_subreddit")
	e.RegisterUser("test_user")

	postID, err := e.PostInSubreddit("test_user", "test_subreddit", "Hello, world!")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if postID != 1 {
		t.Errorf("expected post ID to be 1, got %d", postID)
	}
}

func TestCommentOnPost(t *testing.T) {
	e := engine.NewEngine()

	e.CreateSubreddit("test_subreddit")
	e.RegisterUser("test_user")
	postID, _ := e.PostInSubreddit("test_user", "test_subreddit", "Hello, world!")

	err := e.CommentOnPost("test_user", "test_subreddit", postID, "Nice post!")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestVoting(t *testing.T) {
	e := engine.NewEngine()

	e.CreateSubreddit("test_subreddit")
	e.RegisterUser("test_user")
	postID, _ := e.PostInSubreddit("test_user", "test_subreddit", "Vote on me!")

	err := e.UpvotePost("test_subreddit", postID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = e.DownvotePost("test_subreddit", postID)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestMessaging(t *testing.T) {
	e := engine.NewEngine()

	e.RegisterUser("user1")
	e.RegisterUser("user2")

	err := e.SendMessage("user1", "user2", "Hello!")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	messages, err := e.ListMessages("user2")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(messages))
	}
	if messages[0].Content != "Hello!" {
		t.Errorf("expected message content to be 'Hello!', got '%s'", messages[0].Content)
	}
}
