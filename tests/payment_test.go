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

func setupTestPayment(t *testing.T) uint {
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

func TestProcessPayment(t *testing.T) {
	router := setupTestEnvironment(t)
	playerID := setupTestPayment(t)

	tests := []struct {
		name       string
		payload    map[string]interface{}
		wantStatus int
	}{
		{
			name: "Valid Credit Card Payment",
			payload: map[string]interface{}{
				"player_id": playerID,
				"amount":    100.50,
				"method":    "credit_card",
				"details":   "Test payment",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Valid Bank Transfer",
			payload: map[string]interface{}{
				"player_id": playerID,
				"amount":    200.75,
				"method":    "bank_transfer",
				"details":   "Test bank transfer",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "Invalid Payment Method",
			payload: map[string]interface{}{
				"player_id": playerID,
				"amount":    100.50,
				"method":    "invalid_method",
				"details":   "Should fail",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid Amount",
			payload: map[string]interface{}{
				"player_id": playerID,
				"amount":    -100,
				"method":    "credit_card",
				"details":   "Should fail",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payloadBytes, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/payments", bytes.NewReader(payloadBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("ProcessPayment() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestGetPayment(t *testing.T) {
	router := setupTestEnvironment(t)
	playerID := setupTestPayment(t)

	// Create a test payment
	payment := models.Payment{
		PlayerID: playerID,
		Amount:   100.50,
		Method:   models.PaymentMethodCreditCard,
		Status:   models.PaymentStatusSuccess,
		Details:  "Test payment",
	}
	database.DB.Create(&payment)

	tests := []struct {
		name       string
		paymentID  string
		wantStatus int
	}{
		{
			name:       "Valid Payment",
			paymentID:  fmt.Sprintf("%d", payment.ID),
			wantStatus: http.StatusOK,
		},
		{
			name:       "Invalid Payment ID",
			paymentID:  "999999",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/payments/"+tt.paymentID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("GetPayment() status = %v, want %v", w.Code, tt.wantStatus)
			}
		})
	}
}
