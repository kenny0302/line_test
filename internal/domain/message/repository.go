package message

import (
	"context"
	"main/internal/domain/user"
	"time"
)

// Repository defines the interface for message persistence
type Repository interface {
	// Save persists a message
	Save(ctx context.Context, message *Message) error

	// FindByID retrieves a message by ID
	FindByID(ctx context.Context, id MessageID) (*Message, error)

	// FindByUser retrieves all messages for a specific user
	FindByUser(ctx context.Context, userID user.UserID) ([]*Message, error)

	// FindByTimeRange retrieves messages within a time range
	FindByTimeRange(ctx context.Context, start, end time.Time) ([]*Message, error)

	// Delete removes a message
	Delete(ctx context.Context, id MessageID) error
}
