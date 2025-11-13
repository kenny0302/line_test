package proto

import (
	"encoding/json"
	"testing"
)

func TestUserStruct(t *testing.T) {
	user := User{
		UserId:        "U1234567890",
		DisplayName:   "Test User",
		PictureUrl:    "https://example.com/picture.jpg",
		StatusMessage: "Hello World",
		Language:      "en",
		Message:       "Test message",
	}

	if user.UserId != "U1234567890" {
		t.Errorf("Expected UserId to be U1234567890, got %s", user.UserId)
	}
	if user.DisplayName != "Test User" {
		t.Errorf("Expected DisplayName to be Test User, got %s", user.DisplayName)
	}
	if user.PictureUrl != "https://example.com/picture.jpg" {
		t.Errorf("Expected PictureUrl to be https://example.com/picture.jpg, got %s", user.PictureUrl)
	}
	if user.StatusMessage != "Hello World" {
		t.Errorf("Expected StatusMessage to be Hello World, got %s", user.StatusMessage)
	}
	if user.Language != "en" {
		t.Errorf("Expected Language to be en, got %s", user.Language)
	}
	if user.Message != "Test message" {
		t.Errorf("Expected Message to be Test message, got %s", user.Message)
	}
}

func TestUserJSON(t *testing.T) {
	user := User{
		UserId:      "U1234567890",
		DisplayName: "Test User",
		PictureUrl:  "https://example.com/picture.jpg",
		Message:     "Test message",
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Failed to marshal user: %v", err)
	}

	var unmarshaledUser User
	err = json.Unmarshal(jsonData, &unmarshaledUser)
	if err != nil {
		t.Fatalf("Failed to unmarshal user: %v", err)
	}

	if unmarshaledUser.UserId != user.UserId {
		t.Errorf("Expected UserId to be %s, got %s", user.UserId, unmarshaledUser.UserId)
	}
	if unmarshaledUser.DisplayName != user.DisplayName {
		t.Errorf("Expected DisplayName to be %s, got %s", user.DisplayName, unmarshaledUser.DisplayName)
	}
	if unmarshaledUser.Message != user.Message {
		t.Errorf("Expected Message to be %s, got %s", user.Message, unmarshaledUser.Message)
	}
}

func TestOutputStruct(t *testing.T) {
	output := Output{
		UserId:      "U1234567890",
		DisplayName: "Test User",
		PictureUrl:  "https://example.com/picture.jpg",
	}

	if output.UserId != "U1234567890" {
		t.Errorf("Expected UserId to be U1234567890, got %s", output.UserId)
	}
	if output.DisplayName != "Test User" {
		t.Errorf("Expected DisplayName to be Test User, got %s", output.DisplayName)
	}
	if output.PictureUrl != "https://example.com/picture.jpg" {
		t.Errorf("Expected PictureUrl to be https://example.com/picture.jpg, got %s", output.PictureUrl)
	}
}

func TestOutputJSON(t *testing.T) {
	output := Output{
		UserId:      "U1234567890",
		DisplayName: "Test User",
		PictureUrl:  "https://example.com/picture.jpg",
	}

	jsonData, err := json.Marshal(output)
	if err != nil {
		t.Fatalf("Failed to marshal output: %v", err)
	}

	var unmarshaledOutput Output
	err = json.Unmarshal(jsonData, &unmarshaledOutput)
	if err != nil {
		t.Fatalf("Failed to unmarshal output: %v", err)
	}

	if unmarshaledOutput.UserId != output.UserId {
		t.Errorf("Expected UserId to be %s, got %s", output.UserId, unmarshaledOutput.UserId)
	}
	if unmarshaledOutput.DisplayName != output.DisplayName {
		t.Errorf("Expected DisplayName to be %s, got %s", output.DisplayName, unmarshaledOutput.DisplayName)
	}
	if unmarshaledOutput.PictureUrl != output.PictureUrl {
		t.Errorf("Expected PictureUrl to be %s, got %s", output.PictureUrl, unmarshaledOutput.PictureUrl)
	}
}
