package main

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestSeedDataFromFile(t *testing.T) {
	// Create a temporary seed file
	seedData := SeedData{
		Status: []StatusSeedItem{
			{ID: 0, Name: "SUCCESS"},
			{ID: 1, Name: "FAILED"},
			{ID: 2, Name: "PENDING"},
		},
		Data: []TransactionSeedRow{
			{
				ID:              1001,
				ProductID:       "P001",
				ProductName:     "Product A",
				Amount:          "1500.50",
				CustomerName:    "John Doe",
				Status:          0,
				TransactionDate: "2023-01-15 09:30:00",
				CreateBy:        "admin",
				CreateOn:        "2023-01-15 09:30:00",
			},
			{
				ID:              1002,
				ProductID:       "P002",
				ProductName:     "Product B",
				Amount:          "2500.75",
				CustomerName:    "Jane Smith",
				Status:          1,
				TransactionDate: "2023-01-16 14:45:30",
				CreateBy:        "admin",
				CreateOn:        "2023-01-16 14:45:30",
			},
		},
	}
	
	// Write to temporary file
	tmpFile, err := os.CreateTemp("", "test_seed_*.json")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	
	jsonData, err := json.MarshalIndent(seedData, "", "  ")
	assert.NoError(t, err)
	
	_, err = tmpFile.Write(jsonData)
	assert.NoError(t, err)
	tmpFile.Close()
	
	// Setup test database
	dsn := "host=localhost user=postgres password=postgres dbname=test_seed port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("Failed to connect to test database: %v. Skipping integration test.", err)
	}
	
	// Clean up
	db.Exec("DROP TABLE IF EXISTS transactions")
	db.Exec("DROP TABLE IF EXISTS statuses")
	
	err = db.AutoMigrate(&Status{}, &Transaction{})
	assert.NoError(t, err)
	
	// Temporarily replace global db and seed file
	originalDB := db
	originalSeedFile := os.Getenv("SEED_FILE")
	defer func() { 
		db = originalDB
		os.Setenv("SEED_FILE", originalSeedFile)
	}()
	
	os.Setenv("SEED_FILE", tmpFile.Name())
	
	// Run seed
	seedData()
	
	// Verify data was seeded
	var statusCount int64
	db.Model(&Status{}).Count(&statusCount)
	assert.Equal(t, int64(3), statusCount)
	
	var transactionCount int64
	db.Model(&Transaction{}).Count(&transactionCount)
	assert.Equal(t, int64(2), transactionCount)
	
	// Verify specific records
	var status0 Status
	db.First(&status0, 0)
	assert.Equal(t, 0, status0.ID)
	assert.Equal(t, "SUCCESS", status0.Name)
	
	var status2 Status
	db.First(&status2, 2)
	assert.Equal(t, 2, status2.ID)
	assert.Equal(t, "PENDING", status2.Name)
	
	var transaction1 Transaction
	db.First(&transaction1, 1001)
	assert.Equal(t, int64(1001), transaction1.ID)
	assert.Equal(t, "P001", transaction1.ProductID)
	assert.Equal(t, "Product A", transaction1.ProductName)
	assert.Equal(t, "1500.50", transaction1.AmountText)
	assert.Equal(t, "John Doe", transaction1.CustomerName)
	assert.Equal(t, 0, transaction1.StatusID)
	
	expectedTime1, _ := time.Parse("2006-01-02 15:04:05", "2023-01-15 09:30:00")
	assert.Equal(t, expectedTime1, transaction1.TransactionDate)
	assert.Equal(t, "admin", transaction1.CreateBy)
	assert.Equal(t, expectedTime1, transaction1.CreateOn)
}

func TestSeedData_SkipsIfDataExists(t *testing.T) {
	// Setup test database with existing data
	dsn := "host=localhost user=postgres password=postgres dbname=test_seed_skip port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("Failed to connect to test database: %v. Skipping integration test.", err)
	}
	
	// Clean up and migrate
	db.Exec("DROP TABLE IF EXISTS transactions")
	db.Exec("DROP TABLE IF EXISTS statuses")
	
	err = db.AutoMigrate(&Status{}, &Transaction{})
	assert.NoError(t, err)
	
	// Temporarily replace global db
	originalDB := db
	defer func() { db = originalDB }()
	
	// Insert some existing data
	db.Create(&Status{ID: 99, Name: "EXISTING"})
	db.Create(&Transaction{
		ID:              9999,
		ProductID:       "EXIST",
		ProductName:     "Existing Product",
		AmountText:      "9999",
		CustomerName:    "Existing Customer",
		StatusID:        99,
		TransactionDate: time.Now(),
		CreateBy:        "system",
		CreateOn:        time.Now(),
	})
	
	// Create a seed file that would add different data
	seedData := SeedData{
		Status: []StatusSeedItem{{ID: 100, Name: "NEW"}},
		Data:   []TransactionSeedRow{{ID: 10000, ProductID: "NEW", Status: 100}},
	}
	
	tmpFile, err := os.CreateTemp("", "test_seed_skip_*.json")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	
	jsonData, err := json.Marshal(seedData)
	assert.NoError(t, err)
	
	_, err = tmpFile.Write(jsonData)
	assert.NoError(t, err)
	tmpFile.Close()
	
	originalSeedFile := os.Getenv("SEED_FILE")
	os.Setenv("SEED_FILE", tmpFile.Name())
	defer os.Setenv("SEED_FILE", originalSeedFile)
	
	// Run seed - should skip because data exists
	seedData()
	
	// Verify original data still exists, new data not added
	var statusCount int64
	db.Model(&Status{}).Count(&statusCount)
	assert.Equal(t, int64(1), statusCount) // Only the existing one
	
	var transactionCount int64
	db.Model(&Transaction{}).Count(&transactionCount)
	assert.Equal(t, int64(1), transactionCount) // Only the existing one
	
	var existingStatus Status
	db.First(&existingStatus, 99)
	assert.Equal(t, "EXISTING", existingStatus.Name)
	
	// New status should not exist
	var newStatus Status
	result := db.First(&newStatus, 100)
	assert.Error(t, result.Error)
}
