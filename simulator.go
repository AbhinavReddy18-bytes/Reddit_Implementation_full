package client

import (
	"fmt"
	"log"
	"math/rand"
	"project4/engine"
)

func SimulateClients(e *engine.Engine, userCount, subredditCount, postCount, messageCount int) {

	users := createUsers(e, userCount)

	subreddits := createSubreddits(e, subredditCount)

	GenerateZipfSubreddits(e, subreddits, users)

	simulatePosts(e, users, subreddits, postCount)

	simulateComments(e, users, subreddits)

	simulateMessaging(e, users, messageCount)

	log.Println("Simulation completed successfully.")
}

func createUsers(e *engine.Engine, count int) []string {
	users := []string{}
	for i := 0; i < count; i++ {
		username := fmt.Sprintf("user_%d", i+1)
		err := e.RegisterUser(username)
		if err != nil {
			log.Printf("Error registering user %s: %v", username, err)
		} else {
			connectErr := e.ConnectUser(username)
			if connectErr != nil {
				log.Printf("Error connecting user %s: %v", username, connectErr)
			} else {
				log.Printf("User %s connected successfully.", username)
			}
			users = append(users, username)
		}
	}
	log.Printf("Registered %d users.", len(users))
	return users
}

func createSubreddits(e *engine.Engine, count int) []string {
	subreddits := []string{}
	for i := 0; i < count; i++ {
		name := fmt.Sprintf("subreddit_%d", i+1)
		err := e.CreateSubreddit(name)
		if err != nil {
			log.Printf("Error creating subreddit %s: %v", name, err)
		} else {
			subreddits = append(subreddits, name)
		}
	}
	log.Printf("Created %d subreddits.", len(subreddits))
	return subreddits
}

func simulatePosts(e *engine.Engine, users []string, subreddits []string, count int) {
	for i := 0; i < count; i++ {
		user := users[rand.Intn(len(users))]
		subreddit := subreddits[rand.Intn(len(subreddits))]
		content := fmt.Sprintf("Post #%d by %s in %s", i+1, user, subreddit)
		_, err := e.PostInSubreddit(user, subreddit, content)
		if err != nil {
			log.Printf("Error posting: %v", err)
		}
	}
	log.Printf("Simulated %d posts.", count)
}

func simulateComments(e *engine.Engine, users []string, subreddits []string) {
	for _, subreddit := range subreddits {
		posts, err := e.GetFeed(subreddit, "time", 10)
		if err != nil || len(posts) == 0 {
			continue
		}
		for _, post := range posts {
			user := users[rand.Intn(len(users))]
			commentContent := fmt.Sprintf("Comment by %s on post ID %d", user, post.ID)
			err := e.CommentOnPost(user, subreddit, post.ID, commentContent)
			if err != nil {
				log.Printf("Error commenting: %v", err)
			}
		}
	}
	log.Println("Simulated commenting on posts.")
}

func simulateMessaging(e *engine.Engine, users []string, count int) {
	for i := 0; i < count; i++ {
		sender := users[rand.Intn(len(users))]
		receiver := users[rand.Intn(len(users))]
		if sender == receiver {
			continue
		}
		content := fmt.Sprintf("Message #%d from %s to %s", i+1, sender, receiver)
		err := e.SendMessage(sender, receiver, content)
		if err != nil {
			log.Printf("Error sending message: %v", err)
		}
	}
	log.Printf("Simulated %d direct messages.", count)
}
