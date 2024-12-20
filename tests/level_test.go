package tests

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
)

// We'll use the shared setupTestEnvironment from player_test.go
// so we don't need to import gin directly in this file

func TestCreateLevel(t *testing.T) {
    router := setupTestEnvironment(t)

    tests := []struct {
        name       string
        payload    map[string]interface{}
        wantStatus int
    }{
        {
            name: "Valid Level",
            payload: map[string]interface{}{
                "name": "Beginner",
            },
            wantStatus: http.StatusCreated,
        },
        {
            name: "Invalid Level - Missing Name",
            payload: map[string]interface{}{},
            wantStatus: http.StatusBadRequest,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            payloadBytes, _ := json.Marshal(tt.payload)
            req := httptest.NewRequest("POST", "/levels", bytes.NewReader(payloadBytes))
            req.Header.Set("Content-Type", "application/json")
            
            w := httptest.NewRecorder()
            router.ServeHTTP(w, req)

            if w.Code != tt.wantStatus {
                t.Errorf("CreateLevel() status = %v, want %v", w.Code, tt.wantStatus)
            }
        })
    }
}

func TestListLevels(t *testing.T) {
    router := setupTestEnvironment(t)

    // First create a level
    createPayload := map[string]interface{}{
        "name": "Test Level",
    }
    payloadBytes, _ := json.Marshal(createPayload)
    createReq := httptest.NewRequest("POST", "/levels", bytes.NewReader(payloadBytes))
    createReq.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()
    router.ServeHTTP(w, createReq)

    // Now test listing levels
    req := httptest.NewRequest("GET", "/levels", nil)
    w = httptest.NewRecorder()
    router.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("ListLevels() status = %v, want %v", w.Code, http.StatusOK)
    }

    var response []map[string]interface{}
    if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
        t.Fatalf("Failed to parse response: %v", err)
    }

    if len(response) == 0 {
        t.Error("Expected at least one level in response")
    }
}