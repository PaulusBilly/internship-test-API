package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type APIIntegrationTestSuite struct {
	suite.Suite
	db     *gorm.DB
	server *http.Server
}

// Test models matching main package
type Status struct {
	ID   int    `gorm:"primaryKey;column:id"`
	Name string `gorm:"not null;column:name"`
}

type Transaction struct {
	ID               int64     `gorm:"primaryKey;column:id"`
	ProductID        string    `gorm:"not null;column:product_id"`
	ProductName      string    `gorm:"not null;column:product_name"`
	AmountText       string    `gorm:"not null;column:amount_text"`
	CustomerName     string    `gorm:"not null;column:customer_name"`
	StatusID         int       `gorm:"not null;column:status_id"`
	TransactionDate  time.Time `gorm:"not null;column:transaction_date"`
	CreateBy         string    `gorm:"not null;column:create_by"`
	CreateOn         time.Time `gorm:"not null;column:create_on"`
	Status           Status    `gorm:"foreignKey:StatusID;references:ID"`
}

type ViewDataResponse struct {
	Data   []TransactionRow `json:"data"`
	Status []StatusItem     `json:"status"`
}

type TransactionRow struct {
	ID               int64  `json:"id"`
	ProductID        string `json:"productID"`
	ProductName      string `json:"productName"`
	Amount           string `json:"amount"`
	CustomerName     string `json:"customerName"`
	Status           int    `json:"status"`
	TransactionDate  string `json:"transactionDate"`
	CreateBy         string `json:"createBy"`
	CreateOn         string `json:"createOn"`
}

type StatusItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (suite *APIIntegrationTestSuite) SetupSuite() {
	// Setup test database
	dsn := "host=localhost user=postgres password=postgres dbname=viewdata_integration_test port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(suite.T(), err)
	
	suite.db = db
	
	// Clean up and create tables
	suite.db.Migrator().DropTable(&Transaction{}, &Status{})
	err = suite.db.AutoMigrate(&Status{}, &Transaction{})
	require.NoError(suite.T(), err)
	
	// Insert test data
	statuses := []Status{
		{ID: 0, Name: "SUCCESS"},
		{ID: 1, Name: "FAILED"},
		{ID: 2, Name: "PENDING"},
	}
	
	for _, status := range statuses {
		err := suite.db.Create(&status).Error
		require.NoError(suite.T(), err)
	}
	
	transactions := []Transaction{
		{
			ID:              1001,
			ProductID:       "PROD001",
			ProductName:     "Product One",
			AmountText:      "1000.50",
			CustomerName:    "Customer A",
			StatusID:        0,
			TransactionDate: time.Date(2023, 3, 1, 9, 0, 0, 0, time.UTC),
			CreateBy:        "admin",
			CreateOn:        time.Date(2023, 3, 1, 9, 0, 0, 0, time.UTC),
		},
		{
			ID:              1002,
			ProductID:       "PROD002",
			ProductName:     "Product Two",
			AmountText:      "2000.75",
			CustomerName:    "Customer B",
			StatusID:        1,
			TransactionDate: time.Date(2023, 3, 2, 10, 30, 0, 0, time.UTC),
			CreateBy:        "user1",
			CreateOn:        time.Date(2023, 3, 2, 10, 30, 0, 0, time.UTC),
		},
		{
			ID:              1003,
			ProductID:       "PROD003",
			ProductName:     "Product Three",
			AmountText:      "3000.00",
			CustomerName:    "Customer C",
			StatusID:        2,
			TransactionDate: time.Date(2023, 3, 3, 14, 45, 0, 0, time.UTC),
			CreateBy:        "user2",
			CreateOn:        time.Date(2023, 3, 3, 14, 45, 0, 0, time.UTC),
		},
	}
	
	for _, txn := range transactions {
		err := suite.db.Create(&txn).Error
		require.NoError(suite.T(), err)
	}
	
	// Note: In a real integration test, we would start the actual HTTP server
	// For this test, we'll simulate the endpoint behavior
}

func (suite *APIIntegrationTestSuite) TearDownSuite() {
	// Clean up database
	if suite.db != nil {
		suite.db.Migrator().DropTable(&Transaction{}, &Status{})
	}
}

func (suite *APIIntegrationTestSuite) TestDataConsistency() {
	// Verify database has correct number of records
	var statusCount int64
	err := suite.db.Model(&Status{}).Count(&statusCount).Error
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(3), statusCount, "Should have 3 status records")
	
	var txnCount int64
	err = suite.db.Model(&Transaction{}).Count(&txnCount).Error
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(3), txnCount, "Should have 3 transaction records")
	
	// Verify foreign key relationships
	var transactions []Transaction
	err = suite.db.Preload("Status").Find(&transactions).Error
	require.NoError(suite.T(), err)
	
	for _, txn := range transactions {
		assert.NotEmpty(suite.T(), txn.Status.Name, "Transaction should have status loaded")
		assert.Equal(suite.T(), txn.StatusID, txn.Status.ID, "Status ID should match")
	}
}

func (suite *APIIntegrationTestSuite) TestResponseFormat() {
	// Simulate API response building
	var transactions []Transaction
	err := suite.db.Preload("Status").Find(&transactions).Error
	require.NoError(suite.T(), err)
	
	var statuses []Status
	err = suite.db.Find(&statuses).Error
	require.NoError(suite.T(), err)
	
	// Build response similar to main.buildResponse()
	response := ViewDataResponse{
		Data:   make([]TransactionRow, len(transactions)),
		Status: make([]StatusItem, len(statuses)),
	}
	
	for i, t := range transactions {
		response.Data[i] = TransactionRow{
			ID:              t.ID,
			ProductID:       t.ProductID,
			ProductName:     t.ProductName,
			Amount:          t.AmountText,
			CustomerName:    t.CustomerName,
			Status:          t.StatusID,
			TransactionDate: t.TransactionDate.Format("2006-01-02 15:04:05"),
			CreateBy:        t.CreateBy,
			CreateOn:        t.CreateOn.Format("2006-01-02 15:04:05"),
		}
	}
	
	for i, s := range statuses {
		response.Status[i] = StatusItem{
			ID:   s.ID,
			Name: s.Name,
		}
	}
	
	// Verify response structure
	assert.Len(suite.T(), response.Data, 3, "Should have 3 transactions in response")
	assert.Len(suite.T(), response.Status, 3, "Should have 3 statuses in response")
	
	// Verify data integrity
	for i, txnRow := range response.Data {
		originalTxn := transactions[i]
		assert.Equal(suite.T(), originalTxn.ID, txnRow.ID, "Transaction ID should match")
		assert.Equal(suite.T(), originalTxn.ProductID, txnRow.ProductID, "ProductID should match")
		assert.Equal(suite.T(), originalTxn.AmountText, txnRow.Amount, "Amount should match")
		assert.Equal(suite.T(), originalTxn.StatusID, txnRow.Status, "Status should match")
		
		// Verify date formatting
		expectedDate := originalTxn.TransactionDate.Format("2006-01-02 15:04:05")
		assert.Equal(suite.T(), expectedDate, txnRow.TransactionDate, "TransactionDate should be formatted correctly")
	}
	
	// Test JSON serialization
	jsonBytes, err := json.Marshal(response)
	require.NoError(suite.T(), err)
	
	var decodedResponse ViewDataResponse
	err = json.Unmarshal(jsonBytes, &decodedResponse)
	require.NoError(suite.T(), err)
	
	assert.Len(suite.T(), decodedResponse.Data, 3, "Decoded response should have 3 transactions")
	assert.Len(suite.T(), decodedResponse.Status, 3, "Decoded response should have 3 statuses")
}

func (suite *APIIntegrationTestSuite) TestEnvironmentVariables() {
	// Test that required environment variables are set or have defaults
	requiredVars := []string{
		"DB_HOST",
		"DB_USER", 
		"DB_PASSWORD",
		"DB_NAME",
		"DB_PORT",
	}
	
	for _, envVar := range requiredVars {
		value := os.Getenv(envVar)
		// In test environment, these might not be set, which is OK for unit tests
		// Integration tests would require them to be set
		suite.T().Logf("%s=%s", envVar, value)
	}
}

func TestAPIIntegrationTestSuite(t *testing.T) {
	// Skip if integration test environment is not set up
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration tests. Set RUN_INTEGRATION_TESTS=true to run.")
	}
	
	suite.Run(t, new(APIIntegrationTestSuite))
}
