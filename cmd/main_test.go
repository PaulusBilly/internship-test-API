package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// E1: Empty Database - Happy path seeding
func TestSeedDataEmptyDatabase(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	// Verify database is empty initially
	var count int
	err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM transactions").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)

	// Seed data
	err = seedData(ctx, db)
	require.NoError(t, err)

	// Verify data was seeded
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM transactions").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM statuses").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 2, count)
}

// E2: Existing Data - Seeder should skip when data exists
func TestSeedDataWithExistingData(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	// Insert some data manually first
	_, err := db.ExecContext(ctx, "INSERT INTO statuses (id, name) VALUES (0, 'EXISTING')")
	require.NoError(t, err)

	_, err = db.ExecContext(ctx, `
		INSERT INTO transactions (id, product_id, product_name, amount_text, customer_name,
			status_id, transaction_date, create_by, create_on)
		VALUES (9999, 'test', 'test', '100', 'test', 0, NOW(), 'test', NOW())
	`)
	require.NoError(t, err)

	// Verify initial count
	var count int
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM transactions").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	// Seed should skip due to existing data
	err = seedData(ctx, db)
	require.NoError(t, err)

	// Verify count unchanged
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM transactions").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count)

	// Verify the original data is still there (not overwritten)
	var productName string
	err = db.QueryRowContext(ctx, "SELECT product_name FROM transactions WHERE id = 9999").Scan(&productName)
	require.NoError(t, err)
	assert.Equal(t, "test", productName)
}

// E3: Missing Status Reference - Should error when status not found
func TestSeedDataMissingStatusReference(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	// Create invalid JSON data with missing status reference
	invalidData := `{
		"data": [
			{
				"id": 999,
				"productID": "test",
				"productName": "test",
				"amount": "100",
				"customerName": "test",
				"status": 999,  // Non-existent status
				"transactionDate": "2022-01-01 00:00:00",
				"createBy": "test",
				"createOn": "2022-01-01 00:00:00"
			}
		],
		"status": [
			{"id": 0, "name": "SUCCESS"}
		]
	}`

	// Write temporary JSON file
	tmpFile := "/tmp/test_viewdata.json"
	err := os.WriteFile(tmpFile, []byte(invalidData), 0644)
	require.NoError(t, err)
	defer os.Remove(tmpFile)

	// Test the seeding logic directly
	var count int
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM transactions").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count)

	// This should fail due to foreign key constraint
	jsonData, err := os.ReadFile(tmpFile)
	require.NoError(t, err)

	var viewData ViewDataFile
	err = json.Unmarshal(jsonData, &viewData)
	require.NoError(t, err)

	tx, err := db.Begin()
	require.NoError(t, err)
	defer tx.Rollback()

	// Insert statuses
	for _, s := range viewData.Status {
		_, err := tx.ExecContext(ctx,
			"INSERT INTO statuses (id, name) VALUES ($1, $2) ON CONFLICT (id) DO NOTHING",
			s.ID, s.Name)
		require.NoError(t, err)
	}

	// This should fail due to missing status reference
	for _, r := range viewData.Data {
		txnDate, err := time.Parse("2006-01-02 15:04:05", r.TransactionDate)
		require.NoError(t, err)
		createOn, err := time.Parse("2006-01-02 15:04:05", r.CreateOn)
		require.NoError(t, err)

		_, err = tx.ExecContext(ctx,
			`INSERT INTO transactions (
				id, product_id, product_name, amount_text, customer_name,
				status_id, transaction_date, create_by, create_on
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			r.ID, r.ProductID, r.ProductName, r.Amount, r.CustomerName,
			r.Status, txnDate, r.CreateBy, createOn)
		require.Error(t, err) // Should fail due to foreign key constraint
		assert.Contains(t, err.Error(), "foreign key")
	}
}

// E4: Date Format Mismatch - Should error on invalid date format
func TestSeedDataInvalidDateFormat(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	// Create JSON with invalid date format
	invalidData := `{
		"data": [
			{
				"id": 999,
				"productID": "test",
				"productName": "test",
				"amount": "100",
				"customerName": "test",
				"status": 0,
				"transactionDate": "invalid-date-format",
				"createBy": "test",
				"createOn": "2022-01-01 00:00:00"
			}
		],
		"status": [
			{"id": 0, "name": "SUCCESS"}
		]
	}`

	// Write temporary JSON file
	tmpFile := "/tmp/test_viewdata_invalid.json"
	err := os.WriteFile(tmpFile, []byte(invalidData), 0644)
	require.NoError(t, err)
	defer os.Remove(tmpFile)

	jsonData, err := os.ReadFile(tmpFile)
	require.NoError(t, err)

	var viewData ViewDataFile
	err = json.Unmarshal(jsonData, &viewData)
	require.NoError(t, err)

	// This should fail due to invalid date format
	for _, r := range viewData.Data {
		_, err := time.Parse("2006-01-02 15:04:05", r.TransactionDate)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "parsing time")
	}
}

// E5: Database Connection Failure - Test connection error handling
func TestDatabaseConnectionFailure(t *testing.T) {
	// Test with invalid connection string
	db, err := sql.Open("pgx", "postgres://invalid:invalid@localhost:9999/nonexistent")
	require.NoError(t, err) // Open should succeed

	// Ping should fail
	err = db.Ping()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "connection")
}

// Happy path for BuildResponse
func TestBuildResponseHappyPath(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	// Setup test data manually
	_, err := db.ExecContext(ctx, "INSERT INTO statuses (id, name) VALUES (0, 'SUCCESS'), (1, 'FAILED')")
	require.NoError(t, err)

	testTime := time.Now().Truncate(time.Second)
	_, err = db.ExecContext(ctx, `
		INSERT INTO transactions (id, product_id, product_name, amount_text, customer_name,
			status_id, transaction_date, create_by, create_on)
		VALUES (1001, 'TEST001', 'Test Product', '123.45', 'John Doe', 0, $1, 'system', $2)
	`, testTime, testTime)
	require.NoError(t, err)

	response, err := buildResponse(ctx, db)
	require.NoError(t, err)

	// Verify response structure
	assert.Len(t, response.Data, 1)
	assert.Len(t, response.Status, 2)

	// Verify transaction data
	transaction := response.Data[0]
	assert.Equal(t, int64(1001), transaction.ID)
	assert.Equal(t, "TEST001", transaction.ProductID)
	assert.Equal(t, "Test Product", transaction.ProductName)
	assert.Equal(t, "123.45", transaction.Amount)
	assert.Equal(t, "John Doe", transaction.CustomerName)
	assert.Equal(t, 0, transaction.Status)
	assert.Equal(t, testTime.Format("2006-01-02 15:04:05"), transaction.TransactionDate)
	assert.Equal(t, "system", transaction.CreateBy)
	assert.Equal(t, testTime.Format("2006-01-02 15:04:05"), transaction.CreateOn)

	// Verify status data
	statusMap := make(map[int]string)
	for _, status := range response.Status {
		statusMap[status.ID] = status.Name
	}
	assert.Equal(t, "SUCCESS", statusMap[0])
	assert.Equal(t, "FAILED", statusMap[1])
}

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	testDBURL := os.Getenv("TEST_DATABASE_URL")
	if testDBURL == "" {
		testDBURL = "postgres://postgres:postgres@localhost:5432/viewdata_test?sslmode=disable"
	}

	db, err := sql.Open("pgx", testDBURL)
	require.NoError(t, err)

	// Clean and setup test database
	_, err = db.Exec(`DROP TABLE IF EXISTS transactions CASCADE`)
	require.NoError(t, err)
	_, err = db.Exec(`DROP TABLE IF EXISTS statuses CASCADE`)
	require.NoError(t, err)

	_, err = db.Exec(`
		CREATE TABLE statuses (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL
		)
	`)
	require.NoError(t, err)

	_, err = db.Exec(`
		CREATE TABLE transactions (
			id BIGINT PRIMARY KEY,
			product_id TEXT NOT NULL,
			product_name TEXT NOT NULL,
			amount_text TEXT NOT NULL,
			customer_name TEXT NOT NULL,
			status_id INTEGER NOT NULL REFERENCES statuses(id),
			transaction_date TIMESTAMP NOT NULL,
			create_by TEXT NOT NULL,
			create_on TIMESTAMP NOT NULL
		)
	`)
	require.NoError(t, err)

	return db
}
