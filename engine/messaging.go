package engine

import (
	"fmt"
	"time"
)

func (e *Engine) SendMessage(sender, receiver, content string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	receiverUser, exists := e.Users[receiver]
	if !exists {
		return fmt.Errorf("receiver %s does not exist", receiver)
	}

	message := Message{
		Sender:    sender,
		Receiver:  receiver,
		Content:   content,
		Timestamp: time.Now(),
	}

	receiverUser.Messages = append(receiverUser.Messages, message)
	return nil
}

func (e *Engine) ListMessages(username string) ([]Message, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	user, exists := e.Users[username]
	if !exists {
		return nil, fmt.Errorf("user %s does not exist", username)
	}

	return user.Messages, nil
}

func (e *Engine) ReplyToMessage(sender, receiver, replyContent string) error {
	return e.SendMessage(sender, receiver, replyContent)
}
