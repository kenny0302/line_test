package dto

import "time"

// UserDTO represents a user data transfer object
type UserDTO struct {
	UserID        string    `json:"userid"`
	DisplayName   string    `json:"displayname"`
	PictureURL    string    `json:"pictureurl"`
	StatusMessage string    `json:"statusmessage,omitempty"`
	Language      string    `json:"language,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
}

// MessageDTO represents a message data transfer object
type MessageDTO struct {
	MessageID string    `json:"message_id,omitempty"`
	UserID    string    `json:"userid"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// RegisterUserRequest represents a request to register a new user
type RegisterUserRequest struct {
	UserID        string `json:"userid"`
	DisplayName   string `json:"displayname"`
	PictureURL    string `json:"pictureurl,omitempty"`
	StatusMessage string `json:"statusmessage,omitempty"`
	Language      string `json:"language,omitempty"`
}

// SaveMessageRequest represents a request to save a message
type SaveMessageRequest struct {
	UserID  string `json:"userid"`
	Content string `json:"content"`
}

// PushMessageRequest represents a request to push a message to a user
type PushMessageRequest struct {
	UserID  string `json:"userid"`
	Content string `json:"content"`
}
