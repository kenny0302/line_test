package mongodb

import (
	"context"
	"fmt"
	"main/internal/domain/message"
	"main/internal/domain/user"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MessageRepository implements the message.Repository interface using MongoDB
type MessageRepository struct {
	collection *mongo.Collection
}

// NewMessageRepository creates a new MessageRepository
func NewMessageRepository(db *mongo.Database) *MessageRepository {
	return &MessageRepository{
		collection: db.Collection("message"),
	}
}

// messageDocument represents the MongoDB document structure for a message
type messageDocument struct {
	MessageID string    `bson:"message_id"`
	UserID    string    `bson:"userid"`
	Content   string    `bson:"message"`
	Timestamp string    `bson:"time"` // Stored as string for compatibility
}

// Save persists a message
func (r *MessageRepository) Save(ctx context.Context, m *message.Message) error {
	doc := messageDocument{
		MessageID: m.ID().Value(),
		UserID:    m.UserRef().UserID().Value(),
		Content:   m.Content().Value(),
		Timestamp: fmt.Sprintf("%d", m.Timestamp().Unix()),
	}

	_, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}

	return nil
}

// FindByID retrieves a message by ID
func (r *MessageRepository) FindByID(ctx context.Context, id message.MessageID) (*message.Message, error) {
	var doc messageDocument
	err := r.collection.FindOne(ctx, bson.M{"message_id": id.Value()}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("message not found")
		}
		return nil, fmt.Errorf("failed to find message: %w", err)
	}

	return r.toDomain(doc)
}

// FindByUser retrieves all messages for a specific user
func (r *MessageRepository) FindByUser(ctx context.Context, userID user.UserID) ([]*message.Message, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"userid": userID.Value()})
	if err != nil {
		return nil, fmt.Errorf("failed to find messages: %w", err)
	}
	defer cursor.Close(ctx)

	var messages []*message.Message
	for cursor.Next(ctx) {
		var doc messageDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("failed to decode message: %w", err)
		}

		m, err := r.toDomain(doc)
		if err != nil {
			return nil, err
		}

		messages = append(messages, m)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return messages, nil
}

// FindByTimeRange retrieves messages within a time range
func (r *MessageRepository) FindByTimeRange(ctx context.Context, start, end time.Time) ([]*message.Message, error) {
	startTimestamp := fmt.Sprintf("%d", start.Unix())
	endTimestamp := fmt.Sprintf("%d", end.Unix())

	filter := bson.M{
		"time": bson.M{
			"$gte": startTimestamp,
			"$lte": endTimestamp,
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find messages: %w", err)
	}
	defer cursor.Close(ctx)

	var messages []*message.Message
	for cursor.Next(ctx) {
		var doc messageDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("failed to decode message: %w", err)
		}

		m, err := r.toDomain(doc)
		if err != nil {
			return nil, err
		}

		messages = append(messages, m)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return messages, nil
}

// Delete removes a message
func (r *MessageRepository) Delete(ctx context.Context, id message.MessageID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"message_id": id.Value()})
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}
	return nil
}

// toDomain converts a messageDocument to a domain Message entity
func (r *MessageRepository) toDomain(doc messageDocument) (*message.Message, error) {
	messageID, err := message.NewMessageID(doc.MessageID)
	if err != nil {
		return nil, fmt.Errorf("invalid message ID: %w", err)
	}

	userID, err := user.NewUserID(doc.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	content, err := message.NewContent(doc.Content)
	if err != nil {
		return nil, fmt.Errorf("invalid content: %w", err)
	}

	// Parse timestamp (stored as string)
	var timestamp time.Time
	if doc.Timestamp != "" {
		var unixTime int64
		_, err := fmt.Sscanf(doc.Timestamp, "%d", &unixTime)
		if err == nil {
			timestamp = time.Unix(unixTime, 0)
		} else {
			timestamp = time.Now()
		}
	} else {
		timestamp = time.Now()
	}

	userRef := message.NewUserReference(userID)

	return message.Reconstruct(messageID, userRef, content, timestamp), nil
}
