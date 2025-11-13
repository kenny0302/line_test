package user

import (
	"time"
)

// User represents a LINE bot user (Domain Entity)
type User struct {
	id            UserID
	displayName   DisplayName
	pictureURL    PictureURL
	statusMessage StatusMessage
	language      Language
	createdAt     time.Time
	updatedAt     time.Time
}

// NewUser creates a new User entity
func NewUser(
	id UserID,
	displayName DisplayName,
	pictureURL PictureURL,
	statusMessage StatusMessage,
	language Language,
) *User {
	now := time.Now()
	return &User{
		id:            id,
		displayName:   displayName,
		pictureURL:    pictureURL,
		statusMessage: statusMessage,
		language:      language,
		createdAt:     now,
		updatedAt:     now,
	}
}

// Reconstruct reconstructs a User entity from persistence
func Reconstruct(
	id UserID,
	displayName DisplayName,
	pictureURL PictureURL,
	statusMessage StatusMessage,
	language Language,
	createdAt time.Time,
	updatedAt time.Time,
) *User {
	return &User{
		id:            id,
		displayName:   displayName,
		pictureURL:    pictureURL,
		statusMessage: statusMessage,
		language:      language,
		createdAt:     createdAt,
		updatedAt:     updatedAt,
	}
}

// Getters

// ID returns the user's ID
func (u *User) ID() UserID {
	return u.id
}

// DisplayName returns the user's display name
func (u *User) DisplayName() DisplayName {
	return u.displayName
}

// PictureURL returns the user's picture URL
func (u *User) PictureURL() PictureURL {
	return u.pictureURL
}

// StatusMessage returns the user's status message
func (u *User) StatusMessage() StatusMessage {
	return u.statusMessage
}

// Language returns the user's language
func (u *User) Language() Language {
	return u.language
}

// CreatedAt returns when the user was created
func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

// UpdatedAt returns when the user was last updated
func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

// Business Logic Methods

// UpdateProfile updates the user's profile information
func (u *User) UpdateProfile(displayName DisplayName, pictureURL PictureURL) {
	u.displayName = displayName
	u.pictureURL = pictureURL
	u.updatedAt = time.Now()
}

// UpdateStatusMessage updates the user's status message
func (u *User) UpdateStatusMessage(statusMessage StatusMessage) {
	u.statusMessage = statusMessage
	u.updatedAt = time.Now()
}

// UpdateLanguage updates the user's language preference
func (u *User) UpdateLanguage(language Language) {
	u.language = language
	u.updatedAt = time.Now()
}
