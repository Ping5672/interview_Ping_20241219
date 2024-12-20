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

func setupTestChallenge(t *testing.T) uint {
	// Create a test player with unique name
	player := models.Player{
		Name:  fmt.Sprintf("Test Player %d", time.Now().UnixNano()),
		Level: 1,
	}
	if err := database.DB.Create(&player).Error; err != nil {
		t.Fatalf("Failed to create test player: %v", err)
	}
	return player.ID
}

func TestJoinChallenge(t *testing.T) {
	router := setupTestEnvironment(t)
	playerID := setupTestChallenge(t)

	tests := []struct {
		name       string
		payload    map[string]interface{}
		wantStatus int
	}{
		{
			name: "Valid Challenge",
			payload: map[string]interface{}{
				"player_id": playerID,
				"amount":    20.01,
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "Invalid Amount",
			payload: map[string]interface{}{
				"player_id": playerID,
				"amount":    10.00,
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payloadBytes, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/challenges", bytes.NewReader(payloadBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("JoinChallenge() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestGetChallengeResults(t *testing.T) {
	router := setupTestEnvironment(t)

	req := httptest.NewRequest("GET", "/challenges/results", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("GetChallengeResults() status = %v, want %v", w.Code, http.StatusOK)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if _, ok := response["challenges"]; !ok {
		t.Error("Response missing challenges field")
	}
	if _, ok := response["pool_amount"]; !ok {
		t.Error("Response missing pool_amount field")
	}
}
