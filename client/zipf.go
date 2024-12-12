package client

import (
	"log"
	"math/rand"
	"project4/engine"
	"time"
)

func GenerateZipfSubreddits(e *engine.Engine, subreddits []string, users []string) {
	rand.Seed(time.Now().UnixNano())

	s := 1.1
	v := float64(len(subreddits))
	zipf := rand.NewZipf(rand.New(rand.NewSource(time.Now().UnixNano())), s, v, uint64(len(subreddits)))

	for _, user := range users {
		index := int(zipf.Uint64())
		if index < len(subreddits) {
			err := e.JoinSubreddit(user, subreddits[index])
			if err != nil {
				continue
			}
		}
	}

	for i, subreddit := range subreddits {
		memberCount := len(e.Subreddits[subreddit].Members)
		logSubredditMembership(subreddit, i+1, memberCount)
	}
}

func logSubredditMembership(subreddit string, rank, count int) {
	log.Printf("Subreddit %s (Rank %d): %d members\n", subreddit, rank, count)
}
