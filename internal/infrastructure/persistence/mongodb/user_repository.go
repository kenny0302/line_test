package mongodb

import (
	"context"
	"fmt"
	"main/internal/domain/user"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRepository implements the user.Repository interface using MongoDB
type UserRepository struct {
	collection *mongo.Collection
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("user"),
	}
}

// userDocument represents the MongoDB document structure for a user
type userDocument struct {
	UserID        string    `bson:"userid"`
	DisplayName   string    `bson:"displayname"`
	PictureURL    string    `bson:"pictureurl"`
	StatusMessage string    `bson:"statusmessage,omitempty"`
	Language      string    `bson:"language,omitempty"`
	CreatedAt     time.Time `bson:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at"`
}

// Save persists a user
func (r *UserRepository) Save(ctx context.Context, u *user.User) error {
	doc := userDocument{
		UserID:        u.ID().Value(),
		DisplayName:   u.DisplayName().Value(),
		PictureURL:    u.PictureURL().Value(),
		StatusMessage: u.StatusMessage().Value(),
		Language:      u.Language().Value(),
		CreatedAt:     u.CreatedAt(),
		UpdatedAt:     u.UpdatedAt(),
	}

	filter := bson.M{"userid": doc.UserID}
	update := bson.M{"$set": doc}

	_, err := r.collection.UpdateOne(ctx, filter, update, &mongo.UpdateOptions{
		Upsert: boolPtr(true),
	})

	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	return nil
}

// FindByID retrieves a user by ID
func (r *UserRepository) FindByID(ctx context.Context, id user.UserID) (*user.User, error) {
	var doc userDocument
	err := r.collection.FindOne(ctx, bson.M{"userid": id.Value()}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	return r.toDomain(doc)
}

// FindAll retrieves all users
func (r *UserRepository) FindAll(ctx context.Context) ([]*user.User, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to find users: %w", err)
	}
	defer cursor.Close(ctx)

	var users []*user.User
	for cursor.Next(ctx) {
		var doc userDocument
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("failed to decode user: %w", err)
		}

		u, err := r.toDomain(doc)
		if err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return users, nil
}

// Exists checks if a user exists
func (r *UserRepository) Exists(ctx context.Context, id user.UserID) (bool, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{"userid": id.Value()})
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}
	return count > 0, nil
}

// Delete removes a user
func (r *UserRepository) Delete(ctx context.Context, id user.UserID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"userid": id.Value()})
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// toDomain converts a userDocument to a domain User entity
func (r *UserRepository) toDomain(doc userDocument) (*user.User, error) {
	userID, err := user.NewUserID(doc.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	displayName, err := user.NewDisplayName(doc.DisplayName)
	if err != nil {
		return nil, fmt.Errorf("invalid display name: %w", err)
	}

	pictureURL := user.NewPictureURL(doc.PictureURL)
	statusMessage := user.NewStatusMessage(doc.StatusMessage)
	language := user.NewLanguage(doc.Language)

	return user.Reconstruct(userID, displayName, pictureURL, statusMessage, language, doc.CreatedAt, doc.UpdatedAt), nil
}

func boolPtr(b bool) *bool {
	return &b
}
