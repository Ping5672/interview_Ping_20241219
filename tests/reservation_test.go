package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"interview_Ping_20241219/internal/database"
	"interview_Ping_20241219/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func setupTestReservationData(t *testing.T) (uint, uint) {
	// Create a test player with unique name
	player := models.Player{
		Name:  fmt.Sprintf("Test Player %d", time.Now().UnixNano()),
		Level: 1,
	}
	if err := database.DB.Create(&player).Error; err != nil {
		t.Fatalf("Failed to create test player: %v", err)
	}

	// Create a test room with unique name
	room := models.Room{
		Name:        fmt.Sprintf("Test Room %d", time.Now().UnixNano()),
		Description: "Test Description",
		Status:      "available",
	}
	if err := database.DB.Create(&room).Error; err != nil {
		t.Fatalf("Failed to create test room: %v", err)
	}

	return room.ID, player.ID
}

func TestCreateReservation(t *testing.T) {
	router := setupTestEnvironment(t)

	// Clean up database
	cleanupDatabase()

	roomID, playerID := setupTestReservationData(t)

	tests := []struct {
		name       string
		payload    map[string]interface{}
		wantStatus int
	}{
		{
			name: "Valid Reservation",
			payload: map[string]interface{}{
				"room_id":    roomID,
				"player_id":  playerID,
				"date":       time.Now().Format("2006-01-02"),
				"start_time": "14:00",
				"end_time":   "15:00",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "Invalid Time Format",
			payload: map[string]interface{}{
				"room_id":    roomID,
				"player_id":  playerID,
				"date":       time.Now().Format("2006-01-02"),
				"start_time": "2pm", // Invalid format
				"end_time":   "3pm",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Missing Required Fields",
			payload: map[string]interface{}{
				"room_id": roomID,
				"date":    time.Now().Format("2006-01-02"),
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payloadBytes, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/reservations", bytes.NewReader(payloadBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("CreateReservation() status = %v, want %v, response = %v",
					w.Code, tt.wantStatus, w.Body.String())
			}
		})
	}
}

func TestListReservations(t *testing.T) {
	router := setupTestEnvironment(t)

	// Clean up database
	cleanupDatabase()

	roomID, playerID := setupTestReservationData(t)

	// Create a test reservation
	createPayload := map[string]interface{}{
		"room_id":    roomID,
		"player_id":  playerID,
		"date":       time.Now().Format("2006-01-02"),
		"start_time": "14:00",
		"end_time":   "15:00",
	}

	payloadBytes, _ := json.Marshal(createPayload)
	req := httptest.NewRequest("POST", "/reservations", bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Failed to create test reservation: %v", w.Body.String())
	}

	tests := []struct {
		name       string
		query      string
		wantStatus int
	}{
		{
			name:       "List All Reservations",
			query:      "",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Filter by Room",
			query:      fmt.Sprintf("?room_id=%d", roomID),
			wantStatus: http.StatusOK,
		},
		{
			name:       "Filter by Date",
			query:      fmt.Sprintf("?date=%s", time.Now().Format("2006-01-02")),
			wantStatus: http.StatusOK,
		},
		{
			name:       "With Limit",
			query:      "?limit=5",
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/reservations"+tt.query, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("ListReservations() status = %v, want %v", w.Code, tt.wantStatus)
			}

			var response []map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to parse response: %v", err)
			}
		})
	}
}
