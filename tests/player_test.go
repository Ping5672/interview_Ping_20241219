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

func TestCreatePlayer(t *testing.T) {
	router := setupTestEnvironment(t)

	tests := []struct {
		name       string
		payload    map[string]interface{}
		wantStatus int
	}{
		{
			name: "Valid Player",
			payload: map[string]interface{}{
				"name":  fmt.Sprintf("Test Player %d", time.Now().UnixNano()),
				"level": 1,
			},
			wantStatus: http.StatusCreated,
		},
		{
			name: "Invalid Player - Missing Name",
			payload: map[string]interface{}{
				"level": 1,
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payloadBytes, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/players", bytes.NewReader(payloadBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("CreatePlayer() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}
