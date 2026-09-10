package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"funfillers/backend/config"
	"funfillers/backend/models"
	"funfillers/backend/routes"
)

func TestAddressIsolationBetweenUsers(t *testing.T) {
	cfg := config.Config{
		Port:      "5050",
		JWTSecret: "test_jwt_secret",
		Env:       "test",
	}
	router := routes.SetupRouter(cfg)

	// User A saves address
	addrA := models.Address{
		FullAddress: "123 User A Street, City A",
		City:        "City A",
		State:       "State A",
	}
	bodyA, _ := json.Marshal(addrA)

	reqA := httptest.NewRequest(http.MethodPost, "/api/user/addresses", bytes.NewBuffer(bodyA))
	reqA.Header.Set("Content-Type", "application/json")
	reqA.Header.Set("X-User-ID", "user_a_123")
	recA := httptest.NewRecorder()
	router.ServeHTTP(recA, reqA)

	if recA.Code != http.StatusOK {
		t.Fatalf("Expected status 200 for User A save address, got %d", recA.Code)
	}

	// User B fetches addresses (should be empty for User B)
	reqB := httptest.NewRequest(http.MethodGet, "/api/user/addresses", nil)
	reqB.Header.Set("X-User-ID", "user_b_456")
	recB := httptest.NewRecorder()
	router.ServeHTTP(recB, reqB)

	if recB.Code != http.StatusOK {
		t.Fatalf("Expected status 200 for User B fetch addresses, got %d", recB.Code)
	}

	var listB []models.Address
	if err := json.Unmarshal(recB.Body.Bytes(), &listB); err != nil {
		t.Fatalf("Failed to parse response for User B: %v", err)
	}

	if len(listB) != 0 {
		t.Errorf("Security bug detected! User B received User A's address: %v", listB)
	}

	// User A fetches addresses (should return 1 address)
	reqAGet := httptest.NewRequest(http.MethodGet, "/api/user/addresses", nil)
	reqAGet.Header.Set("X-User-ID", "user_a_123")
	recAGet := httptest.NewRecorder()
	router.ServeHTTP(recAGet, reqAGet)

	var listA []models.Address
	_ = json.Unmarshal(recAGet.Body.Bytes(), &listA)

	if len(listA) != 1 || listA[0].FullAddress != "123 User A Street, City A" {
		t.Errorf("User A did not get their saved address: %v", listA)
	}
}
