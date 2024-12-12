package client

import (
	"fmt"
	"log"
	"project4/engine"
)

type Client struct {
	Username string
	Engine   *engine.Engine
}

func NewClient(username string, engine *engine.Engine) *Client {
	return &Client{
		Username: username,
		Engine:   engine,
	}
}

func (c *Client) Register() error {
	err := c.Engine.RegisterUser(c.Username)
	if err != nil {
		log.Printf("Error registering user %s: %v", c.Username, err)
		return err
	}
	log.Printf("User registered: %s", c.Username)
	return nil
}

func (c *Client) JoinSubreddit(subreddit string) error {
	err := c.Engine.JoinSubreddit(c.Username, subreddit)
	if err != nil {
		log.Printf("Error joining subreddit %s: %v", subreddit, err)
		return err
	}
	log.Printf("%s joined subreddit: %s", c.Username, subreddit)
	return nil
}

func (c *Client) LeaveSubreddit(subreddit string) {
	err := c.Engine.LeaveSubreddit(c.Username, subreddit)
	if err != nil {
		log.Printf("Error leaving subreddit %s: %v\n", subreddit, err)
	} else {
		log.Printf("%s successfully left subreddit %s\n", c.Username, subreddit)
	}
}

func (c *Client) PostInSubreddit(subreddit, content string) (int, error) {
	postID, err := c.Engine.PostInSubreddit(c.Username, subreddit, content)
	if err != nil {
		log.Printf("Error posting in subreddit %s: %v", subreddit, err)
		return 0, err
	}
	log.Printf("Post created by %s in subreddit %s: ID %d", c.Username, subreddit, postID)
	return postID, nil
}

func (c *Client) CommentOnPost(subreddit string, postID int, content string) error {
	err := c.Engine.CommentOnPost(c.Username, subreddit, postID, content)
	if err != nil {
		log.Printf("Error commenting on post ID %d in subreddit %s: %v", postID, subreddit, err)
		return err
	}
	log.Printf("%s commented on post ID %d in subreddit %s", c.Username, postID, subreddit)
	return nil
}

func (c *Client) UpvotePost(subreddit string, postID int) error {
	err := c.Engine.UpvotePost(subreddit, postID)
	if err != nil {
		log.Printf("Error upvoting post ID %d in subreddit %s: %v", postID, subreddit, err)
		return err
	}
	log.Printf("%s upvoted post ID %d in subreddit %s", c.Username, postID, subreddit)
	return nil
}

func (c *Client) DownvotePost(subreddit string, postID int) error {
	err := c.Engine.DownvotePost(subreddit, postID)
	if err != nil {
		log.Printf("Error downvoting post ID %d in subreddit %s: %v", postID, subreddit, err)
		return err
	}
	log.Printf("%s downvoted post ID %d in subreddit %s", c.Username, postID, subreddit)
	return nil
}

func (c *Client) DisplayKarma() {
	karma, err := c.Engine.GetUserKarma(c.Username)
	if err != nil {
		log.Printf("Error fetching karma for %s: %v\n", c.Username, err)
	} else {
		log.Printf("Karma for %s: %d\n", c.Username, karma)
	}
}

func (c *Client) SendMessage(receiver, content string) error {
	err := c.Engine.SendMessage(c.Username, receiver, content)
	if err != nil {
		log.Printf("Error sending message to %s: %v", receiver, err)
		return err
	}
	log.Printf("%s sent a message to %s: %s", c.Username, receiver, content)
	return nil
}

func (c *Client) ListMessages() {
	messages, err := c.Engine.ListMessages(c.Username)
	if err != nil {
		log.Printf("Error listing messages for %s: %v", c.Username, err)
		return
	}
	log.Printf("Messages for %s:", c.Username)
	for _, msg := range messages {
		fmt.Printf("From %s: %s\n", msg.Sender, msg.Content)
	}
}
