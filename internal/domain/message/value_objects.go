package message

import (
	"errors"
	"main/internal/domain/user"
	"strings"
)

// MessageID represents a unique identifier for a message
type MessageID struct {
	value string
}

// NewMessageID creates a new MessageID value object
func NewMessageID(id string) (MessageID, error) {
	if strings.TrimSpace(id) == "" {
		return MessageID{}, errors.New("message ID cannot be empty")
	}
	return MessageID{value: id}, nil
}

// Value returns the string value of the MessageID
func (m MessageID) Value() string {
	return m.value
}

// Content represents the content of a message
type Content struct {
	value string
}

// NewContent creates a new Content value object
func NewContent(content string) (Content, error) {
	if strings.TrimSpace(content) == "" {
		return Content{}, errors.New("message content cannot be empty")
	}
	if len(content) > 5000 {
		return Content{}, errors.New("message content cannot exceed 5000 characters")
	}
	return Content{value: content}, nil
}

// Value returns the string value of the Content
func (c Content) Value() string {
	return c.value
}

// UserReference represents a reference to a user
type UserReference struct {
	userID user.UserID
}

// NewUserReference creates a new UserReference value object
func NewUserReference(userID user.UserID) UserReference {
	return UserReference{userID: userID}
}

// UserID returns the referenced user ID
func (u UserReference) UserID() user.UserID {
	return u.userID
}
