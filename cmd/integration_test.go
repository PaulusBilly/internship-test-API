// +build integration

package main

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationWithRealDatabase(t *testing.T) {
	// Skip if not running integration tests
	if os.Getenv("RUN_INTEGRATION_TESTS") == "" {
		t.Skip("Skipping integration test. Set RUN_INTEGRATION_TESTS=1 to run.")
	}

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=viewdata_test port=5432 sslmode=disable TimeZone=UTC"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Clean up before test
	db.Exec("DROP TABLE IF EXISTS transactions CASCADE")
	db.Exec("DROP TABLE IF EXISTS statuses CASCADE")

	// Auto-migrate
	err = db.AutoMigrate(&Status{}, &Transaction{})
	assert.NoError(t, err)

	// Seed test data
	statuses := []Status{
		{ID: 0, Name: "SUCCESS"},
		{ID: 1, Name: "FAILED"},
	}
	
	for _, status := range statuses {
		assert.NoError(t, db.Create(&status).Error)
	}

	transactions := []Transaction{
		{
			ID:              1001,
			ProductID:       "INT001",
			ProductName:     "Integration Product 1",
			AmountText:      "100.50",
			CustomerName:    "Integration Customer A",
			StatusID:        0,
			TransactionDate: time.Now().Add(-24 * time.Hour),
			CreateBy:        "integration_test",
			CreateOn:        time.Now().Add(-24 * time.Hour),
		},
		{
			ID:              1002,
			ProductID:       "INT002",
			ProductName:     "Integration Product 2",
			AmountText:      "200.75",
			CustomerName:    "Integration Customer B",
			StatusID:        1,
			TransactionDate: time.Now().Add(-12 * time.Hour),
			CreateBy:        "integration_test",
			CreateOn:        time.Now().Add(-12 * time.Hour),
		},
	}
	
	for _, tx := range transactions {
		assert.NoError(t, db.Create(&tx).Error)
	}

	// Test buildResponse with real database
	response := buildResponse(db)
	
	assert.Equal(t, 2, len(response.Data))
	assert.Equal(t, 2, len(response.Status))
	
	// Verify data integrity
	foundProduct1 := false
	foundProduct2 := false
	
	for _, row := range response.Data {
		if row.ProductID == "INT001" {
			foundProduct1 = true
			assert.Equal(t, "Integration Product 1", row.ProductName)
			assert.Equal(t, "100.50", row.Amount)
			assert.Equal(t, 0, row.Status)
		}
		if row.ProductID == "INT002" {
			foundProduct2 = true
			assert.Equal(t, "Integration Product 2", row.ProductName)
			assert.Equal(t, "200.75", row.Amount)
			assert.Equal(t, 1, row.Status)
		}
	}
	
	assert.True(t, foundProduct1, "Product INT001 should be found")
	assert.True(t, foundProduct2, "Product INT002 should be found")
	
	// Clean up after test
	db.Exec("DROP TABLE IF EXISTS transactions CASCADE")
	db.Exec("DROP TABLE IF EXISTS statuses CASCADE")
}

func TestHTTPIntegration(t *testing.T) {
	// Skip if not running integration tests
	if os.Getenv("RUN_INTEGRATION_TESTS") == "" {
		t.Skip("Skipping integration test. Set RUN_INTEGRATION_TESTS=1 to run.")
	}

	// This test would require a running server
	// For now, we'll test the handler logic
	
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=viewdata_test port=5432 sslmode=disable TimeZone=UTC"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Clean up and set up test data
	db.Exec("DROP TABLE IF EXISTS transactions CASCADE")
	db.Exec("DROP TABLE IF EXISTS statuses CASCADE")
	
	err = db.AutoMigrate(&Status{}, &Transaction{})
	assert.NoError(t, err)
	
	// Create minimal test data
	assert.NoError(t, db.Create(&Status{ID: 0, Name: "TEST_STATUS"}).Error)
	assert.NoError(t, db.Create(&Transaction{
		ID:              9999,
		ProductID:       "HTTP_TEST",
		ProductName:     "HTTP Test Product",
		AmountText:      "999",
		CustomerName:    "HTTP Test Customer",
		StatusID:        0,
		TransactionDate: time.Now(),
		CreateBy:        "http_test",
		CreateOn:        time.Now(),
	}).Error)

	// Create test server (in reality, we'd start the actual server)
	// For this test, we'll just verify the handler logic works with real DB
	
	response := buildResponse(db)
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err)
	
	// Verify JSON can be unmarshaled
	var decoded ViewDataResponse
	err = json.Unmarshal(jsonData, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(decoded.Data))
	assert.Equal(t, "HTTP_TEST", decoded.Data[0].ProductID)
	
	// Clean up
	db.Exec("DROP TABLE IF EXISTS transactions CASCADE")
	db.Exec("DROP TABLE IF EXISTS statuses CASCADE")
}

func TestMainFunctionality(t *testing.T) {
	// Test helper functions used in main
	
	// Test time parsing used in seedData
	testTimeStr := "2023-12-25 12:00:00"
	parsedTime, err := time.Parse("2006-01-02 15:04:05", testTimeStr)
	assert.NoError(t, err)
	assert.Equal(t, 2023, parsedTime.Year())
	assert.Equal(t, time.December, parsedTime.Month())
	assert.Equal(t, 25, parsedTime.Day())
	assert.Equal(t, 12, parsedTime.Hour())
	
	// Test JSON structure matches expected API response
	expectedResponse := ViewDataResponse{
		Data: []TransactionRow{
			{
				ID:              1,
				ProductID:       "TEST",
				ProductName:     "Test",
				Amount:          "100",
				CustomerName:    "Test Customer",
				Status:          0,
				TransactionDate: "2023-01-01 00:00:00",
				CreateBy:        "test",
				CreateOn:        "2023-01-01 00:00:00",
			},
		},
		Status: []StatusItem{
			{ID: 0, Name: "SUCCESS"},
			{ID: 1, Name: "FAILED"},
		},
	}
	
	jsonData, err := json.Marshal(expectedResponse)
	assert.NoError(t, err)
	
	// Verify it can be unmarshaled back
	var unmarshaled ViewDataResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(unmarshaled.Data))
	assert.Equal(t, 2, len(unmarshaled.Status))
	assert.Equal(t, "TEST", unmarshaled.Data[0].ProductID)
	assert.Equal(t, "SUCCESS", unmarshaled.Status[0].Name)
}
