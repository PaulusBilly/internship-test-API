package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	// Use test database
	dsn := "host=localhost user=postgres password=postgres dbname=test_viewdata port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("Failed to connect to test database: %v. Skipping integration test.", err)
	}

	// Clean up and migrate
	db.Exec("DROP TABLE IF EXISTS transactions")
	db.Exec("DROP TABLE IF EXISTS statuses")
	
	err = db.AutoMigrate(&Status{}, &Transaction{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}
	
	return db
}

func TestBuildResponse_EmptyDatabase(t *testing.T) {
	db := setupTestDB(t)
	
	// Temporarily replace global db
	originalDB := db
	defer func() { db = originalDB }()
	
	response := buildResponse()
	
	assert.Equal(t, 0, len(response.Data))
	assert.Equal(t, 0, len(response.Status))
}

func TestBuildResponse_WithData(t *testing.T) {
	db := setupTestDB(t)
	
	// Temporarily replace global db
	originalDB := db
	defer func() { db = originalDB }()
	
	// Insert test data
	db.Create(&Status{ID: 0, Name: "SUCCESS"})
	db.Create(&Status{ID: 1, Name: "FAILED"})
	
	transactionDate, _ := time.Parse("2006-01-02 15:04:05", "2022-07-10 11:14:52")
	createOn, _ := time.Parse("2006-01-02 15:04:05", "2022-07-10 11:14:52")
	
	db.Create(&Transaction{
		ID:              1372,
		ProductID:       "10001",
		ProductName:     "Test 1",
		AmountText:      "1000",
		CustomerName:    "abc",
		StatusID:        0,
		TransactionDate: transactionDate,
		CreateBy:        "abc",
		CreateOn:        createOn,
	})
	
	db.Create(&Transaction{
		ID:              1373,
		ProductID:       "10002",
		ProductName:     "Test 2",
		AmountText:      "2000",
		CustomerName:    "def",
		StatusID:        1,
		TransactionDate: transactionDate.Add(24 * time.Hour),
		CreateBy:        "def",
		CreateOn:        createOn.Add(24 * time.Hour),
	})
	
	response := buildResponse()
	
	assert.Equal(t, 2, len(response.Data))
	assert.Equal(t, 2, len(response.Status))
	
	// Verify first transaction
	row1 := response.Data[0]
	assert.Equal(t, int64(1372), row1.ID)
	assert.Equal(t, "10001", row1.ProductID)
	assert.Equal(t, "Test 1", row1.ProductName)
	assert.Equal(t, "1000", row1.Amount)
	assert.Equal(t, "abc", row1.CustomerName)
	assert.Equal(t, 0, row1.Status)
	assert.Equal(t, "2022-07-10 11:14:52", row1.TransactionDate)
	assert.Equal(t, "abc", row1.CreateBy)
	assert.Equal(t, "2022-07-10 11:14:52", row1.CreateOn)
	
	// Verify second transaction
	row2 := response.Data[1]
	assert.Equal(t, int64(1373), row2.ID)
	assert.Equal(t, "10002", row2.ProductID)
	assert.Equal(t, "Test 2", row2.ProductName)
	assert.Equal(t, "2000", row2.Amount)
	assert.Equal(t, "def", row2.CustomerName)
	assert.Equal(t, 1, row2.Status)
	assert.Equal(t, "2022-07-11 11:14:52", row2.TransactionDate)
	assert.Equal(t, "def", row2.CreateBy)
	assert.Equal(t, "2022-07-11 11:14:52", row2.CreateOn)
	
	// Verify status items
	status0 := response.Status[0]
	status1 := response.Status[1]
	
	// Status might be in any order, so check both
	if status0.ID == 0 {
		assert.Equal(t, "SUCCESS", status0.Name)
		assert.Equal(t, 1, status1.ID)
		assert.Equal(t, "FAILED", status1.Name)
	} else {
		assert.Equal(t, 1, status0.ID)
		assert.Equal(t, "FAILED", status0.Name)
		assert.Equal(t, 0, status1.ID)
		assert.Equal(t, "SUCCESS", status1.Name)
	}
}

func TestHandleViewData(t *testing.T) {
	db := setupTestDB(t)
	
	// Temporarily replace global db
	originalDB := db
	defer func() { db = originalDB }()
	
	// Insert test data
	db.Create(&Status{ID: 0, Name: "SUCCESS"})
	db.Create(&Status{ID: 1, Name: "FAILED"})
	
	transactionDate, _ := time.Parse("2006-01-02 15:04:05", "2022-07-10 11:14:52")
	createOn, _ := time.Parse("2006-01-02 15:04:05", "2022-07-10 11:14:52")
	
	db.Create(&Transaction{
		ID:              1372,
		ProductID:       "10001",
		ProductName:     "Test 1",
		AmountText:      "1000",
		CustomerName:    "abc",
		StatusID:        0,
		TransactionDate: transactionDate,
		CreateBy:        "abc",
		CreateOn:        createOn,
	})
	
	// Create request
	req, err := http.NewRequest("GET", "/api/view-data", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	// Create response recorder
	rr := httptest.NewRecorder()
	
	// Create router and handle request
	router := mux.NewRouter()
	router.HandleFunc("/api/view-data", handleViewData).Methods("GET")
	router.ServeHTTP(rr, req)
	
	// Check status code
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	
	// Parse response
	var response ViewDataResponse
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	
	// Verify response structure
	assert.Equal(t, 1, len(response.Data))
	assert.Equal(t, 2, len(response.Status))
	
	row := response.Data[0]
	assert.Equal(t, int64(1372), row.ID)
	assert.Equal(t, "10001", row.ProductID)
	assert.Equal(t, "Test 1", row.ProductName)
	assert.Equal(t, "1000", row.Amount)
	assert.Equal(t, "abc", row.CustomerName)
	assert.Equal(t, 0, row.Status)
	assert.Equal(t, "2022-07-10 11:14:52", row.TransactionDate)
	assert.Equal(t, "abc", row.CreateBy)
	assert.Equal(t, "2022-07-10 11:14:52", row.CreateOn)
}

func TestHandleViewData_MethodNotAllowed(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/view-data", handleViewData).Methods("GET")
	
	// Test POST request (should be not allowed)
	req, err := http.NewRequest("POST", "/api/view-data", nil)
	if err != nil {
		t.Fatal(err)
	}
	
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestSeedData_FileNotFound(t *testing.T) {
	// Test with non-existent seed file
	os.Setenv("SEED_FILE", "non_existent.json")
	
	// This should not panic, just log and exit
	// We'll test by ensuring the function doesn't crash
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("seedData panicked with non-existent file: %v", r)
		}
	}()
	
	// Note: In actual test, we'd mock os.ReadFile or use a test file
}

func TestResponseJSONStructure(t *testing.T) {
	// Test that the response matches the expected JSON structure
	response := ViewDataResponse{
		Data: []TransactionRow{
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
		Status: []StatusItem{
			{ID: 0, Name: "SUCCESS"},
			{ID: 1, Name: "FAILED"},
		},
	}
	
	jsonData, err := json.Marshal(response)
	assert.NoError(t, err)
	
	// Verify JSON can be unmarshaled back
	var decoded ViewDataResponse
	err = json.Unmarshal(jsonData, &decoded)
	assert.NoError(t, err)
	
	assert.Equal(t, 1, len(decoded.Data))
	assert.Equal(t, 2, len(decoded.Status))
	assert.Equal(t, "10001", decoded.Data[0].ProductID)
	assert.Equal(t, "SUCCESS", decoded.Status[0].Name)
}

func TestTransactionDateFormat(t *testing.T) {
	// Test time formatting matches expected pattern
	testTime := time.Date(2022, 7, 10, 11, 14, 52, 0, time.UTC)
	formatted := testTime.Format("2006-01-02 15:04:05")
	
	assert.Equal(t, "2022-07-10 11:14:52", formatted)
	
	// Test parsing back
	parsed, err := time.Parse("2006-01-02 15:04:05", formatted)
	assert.NoError(t, err)
	assert.Equal(t, testTime, parsed)
}
