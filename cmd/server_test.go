package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestBuildResponse(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open in-memory DB: %v", err)
	}

	if err := db.AutoMigrate(&Status{}, &Transaction{}); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	// Seed test data
	if err := db.Create(&Status{ID: 0, Name: "SUCCESS"}).Error; err != nil {
		t.Fatalf("Failed to seed status: %v", err)
	}
	if err := db.Create(&Status{ID: 1, Name: "FAILED"}).Error; err != nil {
		t.Fatalf("Failed to seed status: %v", err)
	}

	if err := db.Create(&Transaction{
		ID:              1372,
		ProductID:       "10001",
		ProductName:     "Test 1",
		AmountText:      "1000",
		CustomerName:    "abc",
		StatusID:        0,
		TransactionDate: mustParseTime("2022-07-10 11:14:52"),
		CreateBy:        "abc",
		CreateOn:        mustParseTime("2022-07-10 11:14:52"),
	}).Error; err != nil {
		t.Fatalf("Failed to seed transaction: %v", err)
	}

	response := buildResponse(db)

	if len(response.Data) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(response.Data))
	}
	if len(response.Status) != 2 {
		t.Errorf("Expected 2 statuses, got %d", len(response.Status))
	}

	// Validate JSON serialization
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	var unmarshaled ViewDataResponse
	if err := json.Unmarshal(jsonBytes, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if unmarshaled.Data[0].ProductID != "10001" {
		t.Errorf("Expected productID '10001', got '%s'", unmarshaled.Data[0].ProductID)
	}
	if unmarshaled.Data[0].Status != 0 {
		t.Errorf("Expected status 0, got %d", unmarshaled.Data[0].Status)
	}
}

func TestAPIEndpoint(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open in-memory DB: %v", err)
	}

	if err := db.AutoMigrate(&Status{}, &Transaction{}); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	if err := db.Create(&Status{ID: 0, Name: "SUCCESS"}).Error; err != nil {
		t.Fatalf("Failed to seed status: %v", err)
	}
	if err := db.Create(&Status{ID: 1, Name: "FAILED"}).Error; err != nil {
		t.Fatalf("Failed to seed status: %v", err)
	}

	if err := db.Create(&Transaction{
		ID:              1372,
		ProductID:       "10001",
		ProductName:     "Test 1",
		AmountText:      "1000",
		CustomerName:    "abc",
		StatusID:        0,
		TransactionDate: mustParseTime("2022-07-10 11:14:52"),
		CreateBy:        "abc",
		CreateOn:        mustParseTime("2022-07-10 11:14:52"),
	}).Error; err != nil {
		t.Fatalf("Failed to seed transaction: %v", err)
	}

	// Override global db for test
	originalDB := db
	defer func() { db = originalDB }()

	req := httptest.NewRequest(http.MethodGet, "/api/view-data", nil)
	w := httptest.NewRecorder()

	http.HandleFunc("/api/view-data", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		response := buildResponse(db)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	http.DefaultServeMux.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var result ViewDataResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 transaction in response, got %d", len(result.Data))
	}
	if result.Data[0].ProductID != "10001" {
		t.Errorf("Expected productID '10001', got '%s'", result.Data[0].ProductID)
	}
	if result.Data[0].Status != 0 {
		t.Errorf("Expected status 0, got %d", result.Data[0].Status)
	}
}

func mustParseTime(s string) time.Time {
	t, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		panic(err)
	}
	return t
}
