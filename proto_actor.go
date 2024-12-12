package engine

import (
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
)

type RegisterUserMessage struct {
	Username string
}

type CreateSubredditMessage struct {
	Name string
}

type JoinSubredditMessage struct {
	Username  string
	Subreddit string
}

type PostMessage struct {
	Username   string
	Subreddit  string
	Content    string
	ResponseCh chan int
}

type LeaveSubredditMessage struct {
	Username  string
	Subreddit string
}

type UpvoteMessage struct {
	Subreddit string
	PostID    int
	Username  string
}

type DownvoteMessage struct {
	Subreddit string
	PostID    int
	Username  string
}

type GetFeedMessage struct {
	Subreddit  string
	SortBy     string
	Limit      int
	ResponseCh chan []Post
}

type EngineActor struct {
	engine *Engine
}

func NewEngineActor() *EngineActor {
	return &EngineActor{
		engine: NewEngine(),
	}
}

func (state *EngineActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *RegisterUserMessage:
		if err := state.engine.RegisterUser(msg.Username); err != nil {
			fmt.Printf("Error registering user: %v\n", err)
		} else {
			fmt.Printf("User registered: %s\n", msg.Username)
		}

	case *CreateSubredditMessage:
		if err := state.engine.CreateSubreddit(msg.Name); err != nil {
			fmt.Printf("Error creating subreddit: %v\n", err)
		} else {
			fmt.Printf("Subreddit created: %s\n", msg.Name)
		}

	case *JoinSubredditMessage:
		if err := state.engine.JoinSubreddit(msg.Username, msg.Subreddit); err != nil {
			fmt.Printf("Error joining subreddit: %v\n", err)
		} else {
			fmt.Printf("User %s joined subreddit %s\n", msg.Username, msg.Subreddit)
		}

	case *LeaveSubredditMessage:
		if err := state.engine.LeaveSubreddit(msg.Username, msg.Subreddit); err != nil {
			fmt.Printf("Error leaving subreddit: %v\n", err)
		} else {
			fmt.Printf("User %s left subreddit %s\n", msg.Username, msg.Subreddit)
		}

	case *UpvoteMessage:
		if err := state.engine.UpvotePost(msg.Subreddit, msg.PostID); err != nil {
			fmt.Printf("Error upvoting post: %v\n", err)
		} else {
			fmt.Printf("User %s upvoted post ID %d in subreddit %s\n", msg.Username, msg.PostID, msg.Subreddit)
			state.engine.UpdateAllUsersKarma()
		}

	case *DownvoteMessage:
		if err := state.engine.DownvotePost(msg.Subreddit, msg.PostID); err != nil {
			fmt.Printf("Error downvoting post: %v\n", err)
		} else {
			fmt.Printf("User %s downvoted post ID %d in subreddit %s\n", msg.Username, msg.PostID, msg.Subreddit)
			state.engine.UpdateAllUsersKarma()
		}

	case *PostMessage:
		postID, err := state.engine.PostInSubreddit(msg.Username, msg.Subreddit, msg.Content)
		if err != nil {
			fmt.Printf("Error posting: %v\n", err)
			msg.ResponseCh <- 0
		} else {
			fmt.Printf("Post created with ID: %d by user %s\n", postID, msg.Username)
			state.engine.UpdateAllUsersKarma()
			msg.ResponseCh <- postID
		}

	case *GetFeedMessage:
		posts, err := state.engine.GetFeed(msg.Subreddit, msg.SortBy, msg.Limit)
		if err != nil {
			fmt.Printf("Error fetching feed: %v\n", err)
			msg.ResponseCh <- nil
		} else {
			msg.ResponseCh <- posts
		}

	default:
		fmt.Printf("Unhandled message: %T\n", msg)
	}
}
