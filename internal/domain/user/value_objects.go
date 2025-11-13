package user

import (
	"errors"
	"strings"
)

// Value Objects

// UserID represents a unique identifier for a user
type UserID struct {
	value string
}

// NewUserID creates a new UserID value object
func NewUserID(id string) (UserID, error) {
	if strings.TrimSpace(id) == "" {
		return UserID{}, errors.New("user ID cannot be empty")
	}
	return UserID{value: id}, nil
}

// Value returns the string value of the UserID
func (u UserID) Value() string {
	return u.value
}

// Equals checks if two UserIDs are equal
func (u UserID) Equals(other UserID) bool {
	return u.value == other.value
}

// DisplayName represents a user's display name
type DisplayName struct {
	value string
}

// NewDisplayName creates a new DisplayName value object
func NewDisplayName(name string) (DisplayName, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return DisplayName{}, errors.New("display name cannot be empty")
	}
	if len(trimmed) > 100 {
		return DisplayName{}, errors.New("display name cannot exceed 100 characters")
	}
	return DisplayName{value: trimmed}, nil
}

// Value returns the string value of the DisplayName
func (d DisplayName) Value() string {
	return d.value
}

// PictureURL represents a user's picture URL
type PictureURL struct {
	value string
}

// NewPictureURL creates a new PictureURL value object
func NewPictureURL(url string) PictureURL {
	return PictureURL{value: url}
}

// Value returns the string value of the PictureURL
func (p PictureURL) Value() string {
	return p.value
}

// IsEmpty checks if the picture URL is empty
func (p PictureURL) IsEmpty() bool {
	return p.value == ""
}

// StatusMessage represents a user's status message
type StatusMessage struct {
	value string
}

// NewStatusMessage creates a new StatusMessage value object
func NewStatusMessage(message string) StatusMessage {
	return StatusMessage{value: message}
}

// Value returns the string value of the StatusMessage
func (s StatusMessage) Value() string {
	return s.value
}

// Language represents a user's language preference
type Language struct {
	value string
}

// NewLanguage creates a new Language value object
func NewLanguage(lang string) Language {
	return Language{value: lang}
}

// Value returns the string value of the Language
func (l Language) Value() string {
	return l.value
}
