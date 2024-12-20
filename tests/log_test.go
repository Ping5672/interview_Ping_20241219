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

func setupTestLog(t *testing.T) uint {
	// Create a test player
	player := models.Player{
		Name:  fmt.Sprintf("Test Player %d", time.Now().UnixNano()),
		Level: 1,
	}
	if err := database.DB.Create(&player).Error; err != nil {
		t.Fatalf("Failed to create test player: %v", err)
	}
	return player.ID
}

func TestCreateLog(t *testing.T) {
	router := setupTestEnvironment(t)
	playerID := setupTestLog(t)

	tests := []struct {
		name       string
		payload    map[string]interface{}
		wantStatus int
	}{
		{
			name: "Valid Log - Register",
			payload: map[string]interface{}{
				"player_id": playerID,
				"action":    "註冊",
				"details":   "New player registration",
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "Invalid Log - Missing Action",
			payload: map[string]interface{}{
				"player_id": playerID,
				"details":   "Missing action",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payloadBytes, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/logs", bytes.NewReader(payloadBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("CreateLog() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestGetLogs(t *testing.T) {
	router := setupTestEnvironment(t)
	playerID := setupTestLog(t)

	// Create some test logs
	testLog := models.GameLog{
		PlayerID: playerID,
		Action:   models.ActionRegister,
		Details:  "Test registration",
	}
	database.DB.Create(&testLog)

	tests := []struct {
		name       string
		query      string
		wantStatus int
	}{
		{
			name:       "Get All Logs",
			query:      "",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Filter by Player",
			query:      fmt.Sprintf("?player_id=%d", playerID),
			wantStatus: http.StatusOK,
		},
		{
			name:       "Filter by Action",
			query:      "?action=註冊",
			wantStatus: http.StatusOK,
		},
		{
			name: "Filter by Time Range",
			query: fmt.Sprintf("?start_time=%s&end_time=%s",
				time.Now().Add(-24*time.Hour).Format(time.RFC3339),
				time.Now().Format(time.RFC3339)),
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/logs"+tt.query, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("GetLogs() status = %v, want %v", w.Code, tt.wantStatus)
			}

			if w.Code == http.StatusOK {
				var response []models.GameLog
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Fatalf("Failed to parse response: %v", err)
				}
			}
		})
	}
}
