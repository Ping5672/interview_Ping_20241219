package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCreateRoom(t *testing.T) {
	router := setupTestEnvironment(t)

	// Clean up database
	cleanupDatabase()

	tests := []struct {
		name       string
		payload    map[string]interface{}
		wantStatus int
	}{
		{
			name: "Valid Room",
			payload: map[string]interface{}{
				"name":        fmt.Sprintf("Test Room %d", time.Now().UnixNano()),
				"description": "A nice gaming room",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "Invalid Room - Missing Name",
			payload: map[string]interface{}{
				"description": "A nice gaming room",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payloadBytes, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/rooms", bytes.NewReader(payloadBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("CreateRoom() status = %v, want %v, response = %v",
					w.Code, tt.wantStatus, w.Body.String())
			}
		})
	}
}

func TestRoomCRUD(t *testing.T) {
	router := setupTestEnvironment(t)

	// Clean up database
	cleanupDatabase()

	// Create a room with a unique name
	createPayload := map[string]interface{}{
		"name":        fmt.Sprintf("Test Room %d", time.Now().UnixNano()),
		"description": "Test Description",
	}
	payloadBytes, _ := json.Marshal(createPayload)
	req := httptest.NewRequest("POST", "/rooms", bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Failed to create test room: %v", w.Body.String())
	}

	var room map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &room)
	roomID := fmt.Sprintf("%.0f", room["id"].(float64))

	// Test Get Room
	t.Run("Get Room", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/rooms/"+roomID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("GetRoom() status = %v, want %v", w.Code, http.StatusOK)
		}
	})

	// Test Update Room
	t.Run("Update Room", func(t *testing.T) {
		updatePayload := map[string]interface{}{
			"name":        "Updated Room",
			"description": "Updated Description",
		}
		payloadBytes, _ := json.Marshal(updatePayload)
		req := httptest.NewRequest("PUT", "/rooms/"+roomID, bytes.NewReader(payloadBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("UpdateRoom() status = %v, want %v", w.Code, http.StatusOK)
		}
	})

	// Test Delete Room
	t.Run("Delete Room", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/rooms/"+roomID, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("DeleteRoom() status = %v, want %v", w.Code, http.StatusOK)
		}
	})
}
