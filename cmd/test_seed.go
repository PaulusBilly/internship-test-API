package main

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSeedData(t *testing.T) {
	// Setup test database
	os.Setenv("DATABASE_URL", "host=localhost user=test password=test dbname=test_seed_data port=5432 sslmode=disable TimeZone=UTC")
	
	// Create test data
	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	
	err = db.AutoMigrate(&Status{}, &Transaction{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}
	
	// Create seed data
	seedData := SeedData{
		Status: []StatusSeedItem{
			{ID: 0, Name: "SUCCESS"},
			{ID: 1, Name: "FAILED"},
		},
		Data: []TransactionSeedRow{
			{
				ID:              1372,
				ProductID:       "10001",
				ProductName:     "Test 1",
				Amount:          "1000",
				CustomerName:    "abc",
				Status:          0,
				TransactionDate: "2022-07-10 11:14:52",
				CreateBy:        "abc",
				CreateOn:        "2022-07-10 11:14:52",
			},
		},
	}
	
	// Save to JSON file
	jsonData, err := json.Marshal(seedData)
	if err != nil {
		t.Fatalf("Failed to marshal seed data: %v", err)
	}
	
	err = os.WriteFile("test_seed.json", jsonData, 0644)
	if err != nil {
		t.Fatalf("Failed to write seed file: %v", err)
	}
	defer os.Remove("test_seed.json")
	
	// Set environment variable for seed file
	os.Setenv("SEED_FILE", "test_seed.json")
	
	// Seed data
	seedData()
	
	// Verify data was seeded
	var count int64
	db.Model(&Transaction{}).Count(&count)
	assert.Equal(t, int64(1), count)
	
	var transaction Transaction
	db.First(&transaction, 1372)
	assert.Equal(t, int64(1372), transaction.ID)
	assert.Equal(t, "10001", transaction.ProductID)
	assert.Equal(t, "Test 1", transaction.ProductName)
	assert.Equal(t, "1000", transaction.AmountText)
	assert.Equal(t, "abc", transaction.CustomerName)
	assert.Equal(t, 0, transaction.StatusID)
	assert.Equal(t, time.Date(2022, 7, 10, 11, 14, 52, 0, time.UTC), transaction.TransactionDate)
	assert.Equal(t, "abc", transaction.CreateBy)
	assert.Equal(t, time.Date(2022, 7, 10, 11, 14, 52, 0, time.UTC), transaction.CreateOn)
	
	// Verify status
	var status Status
	db.First(&status, 0)
	assert.Equal(t, 0, status.ID)
	assert.Equal(t, "SUCCESS", status.Name)
}
