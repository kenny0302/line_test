package message

import (
	"time"
)

// Message represents a chat message (Domain Entity)
type Message struct {
	id        MessageID
	userRef   UserReference
	content   Content
	timestamp time.Time
}

// NewMessage creates a new Message entity
func NewMessage(
	id MessageID,
	userRef UserReference,
	content Content,
) *Message {
	return &Message{
		id:        id,
		userRef:   userRef,
		content:   content,
		timestamp: time.Now(),
	}
}

// Reconstruct reconstructs a Message entity from persistence
func Reconstruct(
	id MessageID,
	userRef UserReference,
	content Content,
	timestamp time.Time,
) *Message {
	return &Message{
		id:        id,
		userRef:   userRef,
		content:   content,
		timestamp: timestamp,
	}
}

// Getters

// ID returns the message's ID
func (m *Message) ID() MessageID {
	return m.id
}

// UserRef returns the reference to the user who sent the message
func (m *Message) UserRef() UserReference {
	return m.userRef
}

// Content returns the message content
func (m *Message) Content() Content {
	return m.content
}

// Timestamp returns when the message was created
func (m *Message) Timestamp() time.Time {
	return m.timestamp
}

// Business Logic Methods

// IsSentBy checks if the message was sent by a specific user
func (m *Message) IsSentBy(userRef UserReference) bool {
	return m.userRef.UserID().Equals(userRef.UserID())
}

// IsOlderThan checks if the message is older than a given duration
func (m *Message) IsOlderThan(duration time.Duration) bool {
	return time.Since(m.timestamp) > duration
}
