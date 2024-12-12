package engine

import (
	"fmt"
	"project4/performance"
	"sort"
	"sync"
	"time"
)

type User struct {
	Username  string
	Messages  []Message
	Karma     int
	Connected bool
}

type Subreddit struct {
	Name    string
	Members map[string]*User
	Posts   []Post
}

type Post struct {
	ID        int
	Author    string
	Content   string
	Comments  []*Comment
	Upvotes   int
	Downvotes int
	Timestamp time.Time
}

type Comment struct {
	ID        int
	Author    string
	Content   string
	Replies   []*Comment
	Timestamp time.Time
}

type Message struct {
	Sender    string
	Receiver  string
	Content   string
	Timestamp time.Time
}

type Engine struct {
	mu           sync.Mutex
	Users        map[string]*User
	Subreddits   map[string]*Subreddit
	PostCount    int
	CommentCount int
	metrics      *performance.Metrics
}

func NewEngine() *Engine {
	return &Engine{
		Users:      make(map[string]*User),
		Subreddits: make(map[string]*Subreddit),
	}
}

func (e *Engine) SetMetrics(metrics *performance.Metrics) {
	e.metrics = metrics
}

func (e *Engine) RegisterUser(username string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, exists := e.Users[username]; exists {
		return fmt.Errorf("user %s already exists", username)
	}

	e.Users[username] = &User{Username: username, Karma: 0}
	e.metrics.IncrementOperation()
	return nil
}

func (e *Engine) CreateSubreddit(name string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, exists := e.Subreddits[name]; exists {
		return fmt.Errorf("subreddit %s already exists", name)
	}

	e.Subreddits[name] = &Subreddit{Name: name, Members: make(map[string]*User)}
	e.metrics.IncrementOperation()
	return nil
}

func (e *Engine) JoinSubreddit(username, subreddit string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	user, userExists := e.Users[username]
	sub, subExists := e.Subreddits[subreddit]
	if !userExists || !subExists {
		return fmt.Errorf("invalid user or subreddit")
	}

	sub.Members[username] = user
	e.metrics.IncrementOperation()
	return nil
}

func (e *Engine) LeaveSubreddit(username, subreddit string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	sub, exists := e.Subreddits[subreddit]
	if !exists {
		return fmt.Errorf("subreddit does not exist")
	}

	_, memberExists := sub.Members[username]
	if !memberExists {
		return fmt.Errorf("user %s is not a member of subreddit %s", username, subreddit)
	}

	delete(sub.Members, username)
	e.metrics.IncrementOperation()
	return nil
}

func (e *Engine) PostInSubreddit(username, subreddit, content string) (int, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	user, userExists := e.Users[username]
	sub, subExists := e.Subreddits[subreddit]
	if !userExists || !subExists {
		return 0, fmt.Errorf("subreddit does not exist")
	}

	if !user.Connected {
		return 0, fmt.Errorf("user %s is not connected", username)
	}

	post := Post{
		ID:        e.PostCount + 1,
		Author:    username,
		Content:   content,
		Comments:  []*Comment{},
		Timestamp: time.Now(),
	}
	e.PostCount++
	sub.Posts = append(sub.Posts, post)
	e.metrics.IncrementOperation()
	return post.ID, nil
}

func (e *Engine) CommentOnPost(username, subreddit string, postID int, content string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	user, userExists := e.Users[username]
	sub, subExists := e.Subreddits[subreddit]
	if !userExists || !subExists {
		return fmt.Errorf("subreddit does not exist")
	}

	if !user.Connected {
		return fmt.Errorf("user %s is not connected", username)
	}

	for i := range sub.Posts {
		if sub.Posts[i].ID == postID {
			comment := &Comment{
				ID:        e.CommentCount + 1,
				Author:    username,
				Content:   content,
				Replies:   []*Comment{},
				Timestamp: time.Now(),
			}
			sub.Posts[i].Comments = append(sub.Posts[i].Comments, comment)
			e.CommentCount++
			e.metrics.IncrementOperation()
			fmt.Printf("Comment added by user %s on post %d in subreddit %s\n", username, postID, subreddit)
			return nil
		}
	}

	return fmt.Errorf("post not found")
}

func (e *Engine) ReplyToComment(subreddit string, postID, parentCommentID int, username, content string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	sub, exists := e.Subreddits[subreddit]
	if !exists {
		return fmt.Errorf("subreddit does not exist")
	}

	for i := range sub.Posts {
		if sub.Posts[i].ID == postID {
			parent := findCommentByID(sub.Posts[i].Comments, parentCommentID)
			if parent == nil {
				return fmt.Errorf("comment not found")
			}
			reply := &Comment{
				ID:        e.CommentCount + 1,
				Author:    username,
				Content:   content,
				Timestamp: time.Now(),
			}
			parent.Replies = append(parent.Replies, reply)
			e.CommentCount++
			e.metrics.IncrementOperation()
			return nil
		}
	}

	return fmt.Errorf("post not found")
}

func findCommentByID(comments []*Comment, id int) *Comment {
	for _, comment := range comments {
		if comment.ID == id {
			return comment
		}
		for _, reply := range comment.Replies {
			if found := findCommentByID(reply.Replies, id); found != nil {
				return found
			}
		}
	}
	return nil
}

func (e *Engine) UpvotePost(subreddit string, postID int) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	sub, exists := e.Subreddits[subreddit]
	if !exists {
		return fmt.Errorf("subreddit does not exist")
	}

	for i := range sub.Posts {
		if sub.Posts[i].ID == postID {
			sub.Posts[i].Upvotes++
			e.metrics.IncrementOperation()
			fmt.Printf("Post %d in subreddit %s upvoted. Total upvotes: %d\n", postID, subreddit, sub.Posts[i].Upvotes)
			return nil
		}
	}
	return fmt.Errorf("post not found")
}

func (e *Engine) DownvotePost(subreddit string, postID int) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	sub, exists := e.Subreddits[subreddit]
	if !exists {
		return fmt.Errorf("subreddit does not exist")
	}

	for i := range sub.Posts {
		if sub.Posts[i].ID == postID {
			sub.Posts[i].Downvotes++
			e.metrics.IncrementOperation()
			fmt.Printf("Post %d in subreddit %s downvoted. Total downvotes: %d\n", postID, subreddit, sub.Posts[i].Downvotes)
			return nil
		}
	}
	return fmt.Errorf("post not found")
}

func (e *Engine) ComputeKarma(username string) int {
	e.mu.Lock()
	defer e.mu.Unlock()

	user, exists := e.Users[username]
	if !exists {
		return 0
	}

	karma := 0
	for _, sub := range e.Subreddits {
		for _, post := range sub.Posts {
			if post.Author == username {
				karma += post.Upvotes - post.Downvotes
			}
			for _, comment := range post.Comments {
				karma += computeCommentKarma(comment, username)
			}
		}
	}

	user.Karma = karma
	fmt.Printf("Computed karma for user %s: %d\n", username, karma)
	e.metrics.IncrementOperation()
	return karma
}

func computeCommentKarma(comment *Comment, username string) int {
	karma := 0
	if comment.Author == username {
		karma++
	}
	for _, reply := range comment.Replies {
		karma += computeCommentKarma(reply, username)
	}
	return karma
}

func (e *Engine) UpdateAllUsersKarma() {
	e.mu.Lock()
	defer e.mu.Unlock()

	for username := range e.Users {
		e.ComputeKarma(username)
	}
}

func (e *Engine) GetUserKarma(username string) (int, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	user, exists := e.Users[username]
	if !exists {
		return 0, fmt.Errorf("user %s does not exist", username)
	}

	return user.Karma, nil
}

func (e *Engine) GetFeed(subreddit string, sortBy string, limit int) ([]Post, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	sub, exists := e.Subreddits[subreddit]
	if !exists {
		return nil, fmt.Errorf("subreddit not found")
	}

	posts := append([]Post{}, sub.Posts...)
	switch sortBy {
	case "upvotes":
		sort.Slice(posts, func(i, j int) bool { return posts[i].Upvotes > posts[j].Upvotes })
	case "time":
		sort.Slice(posts, func(i, j int) bool { return posts[i].Timestamp.After(posts[j].Timestamp) })
	default:
		return nil, fmt.Errorf("invalid sort criteria")
	}

	if len(posts) > limit {
		posts = posts[:limit]
	}
	e.metrics.IncrementOperation()
	return posts, nil
}

func (e *Engine) ConnectUser(username string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	user, exists := e.Users[username]
	if !exists {
		return fmt.Errorf("user %s does not exist", username)
	}

	user.Connected = true
	fmt.Printf("User %s is now connected.\n", username)
	return nil
}

func (e *Engine) DisconnectUser(username string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	user, exists := e.Users[username]
	if !exists {
		return fmt.Errorf("user %s does not exist", username)
	}

	user.Connected = false
	fmt.Printf("User %s is now disconnected.\n", username)
	return nil
}
