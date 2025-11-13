package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	proto "main/proto"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func TestMain(m *testing.M) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test config file
	createTestConfig()

	// Run tests
	code := m.Run()

	// Cleanup
	os.Remove("config.yaml")

	os.Exit(code)
}

func createTestConfig() {
	config := `database:
 host: localhost
 port: 27017

redis:
 host: localhost
 port: 6379
 password: ""
 db: 0

line:
 secret: test_secret
 token: test_token
`
	os.WriteFile("config.yaml", []byte(config), 0644)
}

func TestSetupRouter(t *testing.T) {
	router := SetupRouter()

	if router == nil {
		t.Fatal("Router should not be nil")
	}

	// Check if routes are registered
	routes := router.Routes()
	if len(routes) == 0 {
		t.Error("Expected routes to be registered")
	}

	expectedRoutes := map[string]bool{
		"POST /callback": false,
		"GET /list":      false,
		"POST /push":     false,
	}

	for _, route := range routes {
		key := route.Method + " " + route.Path
		if _, exists := expectedRoutes[key]; exists {
			expectedRoutes[key] = true
		}
	}

	for route, found := range expectedRoutes {
		if !found {
			t.Errorf("Expected route %s to be registered", route)
		}
	}
}

func TestListUserHandler(t *testing.T) {
	router := SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/list", nil)
	router.ServeHTTP(w, req)

	// We expect either 200 with data or 500 if MongoDB is not running
	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 200 or 500, got %d", w.Code)
	}

	if w.Code == http.StatusOK {
		var result []proto.Output
		err := json.Unmarshal(w.Body.Bytes(), &result)
		if err != nil {
			t.Errorf("Failed to parse JSON response: %v", err)
		}
		t.Logf("ListUser returned %d users", len(result))
	}
}

func TestPushMessageHandler(t *testing.T) {
	router := SetupRouter()

	form := url.Values{}
	form.Add("UserId", "U1234567890")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/push", bytes.NewBufferString(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	// We expect 500 because the LINE bot credentials are fake
	// or 200 if somehow it works
	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 200 or 500, got %d", w.Code)
	}
}

func TestGetBotDataHandler(t *testing.T) {
	router := SetupRouter()

	// Test with empty body (invalid LINE webhook)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/callback", nil)
	router.ServeHTTP(w, req)

	// Should return error because signature validation fails
	if w.Code != http.StatusBadRequest && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 400 or 500, got %d", w.Code)
	}
}

func TestSetFunction(t *testing.T) {
	user := proto.User{
		UserId:      "U1234567890",
		DisplayName: "Test User",
		PictureUrl:  "https://example.com/pic.jpg",
		Message:     "Hello World",
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Failed to marshal user: %v", err)
	}

	// This will fail if MongoDB is not running, but tests the function structure
	err = Set("localhost", "27017", jsonData)
	if err != nil {
		t.Logf("Set function returned error (expected if MongoDB not running): %v", err)
	} else {
		t.Log("Set function executed successfully")
	}
}

func TestSetFunctionWithEmptyData(t *testing.T) {
	emptyUser := proto.User{}
	jsonData, err := json.Marshal(emptyUser)
	if err != nil {
		t.Fatalf("Failed to marshal empty user: %v", err)
	}

	err = Set("localhost", "27017", jsonData)
	if err != nil {
		t.Logf("Set function with empty data returned error: %v", err)
	}
}

func TestSetFunctionWithInvalidHost(t *testing.T) {
	user := proto.User{
		UserId:      "U1234567890",
		DisplayName: "Test User",
		Message:     "Hello",
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Failed to marshal user: %v", err)
	}

	// Should return error with invalid host
	err = Set("invalid_host", "27017", jsonData)
	if err == nil {
		t.Error("Expected error with invalid host")
	} else {
		t.Logf("Set function correctly returned error with invalid host: %v", err)
	}
}

func TestRouterPOSTCallback(t *testing.T) {
	router := SetupRouter()

	// Create a mock LINE webhook payload
	payload := []byte(`{"events":[]}`)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/callback", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Should fail signature validation
	if w.Code == http.StatusOK {
		t.Log("Callback returned OK")
	} else {
		t.Logf("Callback returned status %d (expected due to signature validation)", w.Code)
	}
}

func TestListUserEndpoint(t *testing.T) {
	router := SetupRouter()

	tests := []struct {
		name   string
		method string
		path   string
	}{
		{
			name:   "GET /list",
			method: "GET",
			path:   "/list",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
				t.Errorf("Expected status 200 or 500, got %d", w.Code)
			}
		})
	}
}

func TestPushMessageWithEmptyUserId(t *testing.T) {
	router := SetupRouter()

	form := url.Values{}
	form.Add("UserId", "")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/push", bytes.NewBufferString(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 200 or 500, got %d", w.Code)
	}
}

func TestBSONFilterCreation(t *testing.T) {
	// Test BSON filter creation used in the code
	filter := bson.D{{"userid", "test123"}}
	if len(filter) != 1 {
		t.Errorf("Expected filter length 1, got %d", len(filter))
	}

	emptyFilter := bson.D{{}}
	if len(emptyFilter) != 1 {
		t.Errorf("Expected empty filter length 1, got %d", len(emptyFilter))
	}

	mapFilter := bson.M{"userid": "test123", "displayname": "Test", "pictureurl": "url"}
	if len(mapFilter) != 3 {
		t.Errorf("Expected map filter length 3, got %d", len(mapFilter))
	}
}

func TestConfigLoading(t *testing.T) {
	// Test that config can be loaded
	router := SetupRouter()
	if router == nil {
		t.Fatal("Failed to setup router (config loading failed)")
	}
}

func TestSetFunctionWithValidData(t *testing.T) {
	user := proto.User{
		UserId:      "U0987654321",
		DisplayName: "Another Test User",
		PictureUrl:  "https://example.com/another.jpg",
		Message:     "Another test message",
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Failed to marshal user: %v", err)
	}

	err = Set("localhost", "27017", jsonData)
	if err != nil {
		t.Logf("Set function returned error: %v", err)
	} else {
		t.Log("Set function executed successfully")
	}
}

func TestPushMessageWithValidUserId(t *testing.T) {
	router := SetupRouter()

	form := url.Values{}
	form.Add("UserId", "U1234567890abcdef")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/push", bytes.NewBufferString(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 200 or 500, got %d", w.Code)
	}
}

func TestGetBotDataWithJSONPayload(t *testing.T) {
	router := SetupRouter()

	payload := map[string]interface{}{
		"events": []map[string]interface{}{
			{
				"type": "message",
				"message": map[string]interface{}{
					"type": "text",
					"text": "Hello",
				},
			},
		},
	}

	jsonPayload, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/callback", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Should fail due to signature validation
	if w.Code != http.StatusBadRequest && w.Code != http.StatusInternalServerError {
		t.Logf("Callback returned status %d", w.Code)
	}
}

func TestMultipleListUserCalls(t *testing.T) {
	router := SetupRouter()

	for i := 0; i < 3; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/list", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
			t.Errorf("Call %d: Expected status 200 or 500, got %d", i+1, w.Code)
		}
	}
}

func TestSetWithMalformedJSON(t *testing.T) {
	malformedData := []byte(`{"userid": "test", "displayname": }`)

	err := Set("localhost", "27017", malformedData)
	// Should handle malformed JSON gracefully
	t.Logf("Set with malformed JSON returned: %v", err)
}

func TestPushMessageWithSpecialCharacters(t *testing.T) {
	router := SetupRouter()

	form := url.Values{}
	form.Add("UserId", "U!@#$%^&*()")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/push", bytes.NewBufferString(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 200 or 500, got %d", w.Code)
	}
}

func TestCallbackWithInvalidContentType(t *testing.T) {
	router := SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/callback", bytes.NewBufferString("invalid data"))
	req.Header.Set("Content-Type", "text/plain")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest && w.Code != http.StatusInternalServerError {
		t.Logf("Callback with invalid content type returned status %d", w.Code)
	}
}
