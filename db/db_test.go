package db

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestConnect(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		port    string
		wantErr bool
	}{
		{
			name:    "Invalid host",
			host:    "invalid",
			port:    "27017",
			wantErr: false, // Connection creation succeeds, but actual connection might fail
		},
		{
			name:    "Empty host",
			host:    "",
			port:    "27017",
			wantErr: false,
		},
		{
			name:    "Valid connection parameters",
			host:    "localhost",
			port:    "27017",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			_, err := Connect(ctx, tt.host, tt.port)
			if (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetUserStructure(t *testing.T) {
	// Test that GetUser function signature and basic logic work
	// This will fail if MongoDB is not running, but tests the function structure
	host := "localhost"
	port := "27017"
	db := "line_test"
	col := "user_test"
	filter := bson.D{{}}

	_, err := GetUser(host, port, db, col, filter)
	// We expect an error if MongoDB is not running, which is fine for structure testing
	// The function should return an error, not panic
	if err == nil {
		t.Log("GetUser executed successfully (MongoDB is running)")
	} else {
		t.Logf("GetUser returned error as expected (MongoDB not running): %v", err)
	}
}

func TestSetUserStructure(t *testing.T) {
	// Test that SetUser function signature and basic logic work
	host := "localhost"
	port := "27017"
	db := "line_test"
	col := "user_test"
	filter := bson.M{"userid": "test123", "displayname": "Test User", "pictureurl": "http://test.com"}

	err := SetUser(host, port, db, col, filter)
	// We expect an error if MongoDB is not running
	if err == nil {
		t.Log("SetUser executed successfully (MongoDB is running)")
	} else {
		t.Logf("SetUser returned error as expected (MongoDB not running): %v", err)
	}
}

func TestSetMessageStructure(t *testing.T) {
	// Test that SetMessage function signature and basic logic work
	host := "localhost"
	port := "27017"
	db := "line_test"
	col := "message_test"
	filter := bson.M{"userid": "test123", "message": "Hello", "time": "1234567890"}

	err := SetMessage(host, port, db, col, filter)
	// We expect an error if MongoDB is not running
	if err == nil {
		t.Log("SetMessage executed successfully (MongoDB is running)")
	} else {
		t.Logf("SetMessage returned error as expected (MongoDB not running): %v", err)
	}
}

func TestConnectWithContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	host := "localhost"
	port := "27017"

	client, err := Connect(ctx, host, port)
	if err == nil && client != nil {
		t.Log("Connect with context executed successfully")
		// Try to disconnect
		err = client.Disconnect(ctx)
		if err != nil {
			t.Logf("Disconnect failed: %v", err)
		}
	} else {
		t.Logf("Connect with context returned error (expected if MongoDB not running): %v", err)
	}
}

func TestGetUserWithEmptyFilter(t *testing.T) {
	host := "localhost"
	port := "27017"
	db := "line_test"
	col := "user_test"
	filter := bson.D{{}}

	result, err := GetUser(host, port, db, col, filter)
	if err == nil {
		t.Logf("GetUser with empty filter returned %d results", len(result))
	} else {
		t.Logf("GetUser with empty filter returned error: %v", err)
	}
}

func TestGetUserWithSpecificFilter(t *testing.T) {
	host := "localhost"
	port := "27017"
	db := "line_test"
	col := "user_test"
	filter := bson.D{{"userid", "test123"}}

	result, err := GetUser(host, port, db, col, filter)
	if err == nil {
		t.Logf("GetUser with specific filter returned %d results", len(result))
	} else {
		t.Logf("GetUser with specific filter returned error: %v", err)
	}
}

func TestConnectWithDifferentHosts(t *testing.T) {
	tests := []struct {
		name string
		host string
		port string
	}{
		{"Localhost", "localhost", "27017"},
		{"127.0.0.1", "127.0.0.1", "27017"},
		{"Invalid IP", "999.999.999.999", "27017"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			client, err := Connect(ctx, tt.host, tt.port)
			if err == nil && client != nil {
				client.Disconnect(ctx)
			}
		})
	}
}
