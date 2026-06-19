package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"internship-test-api/internal/db"
	"internship-test-api/internal/handler"
	"internship-test-api/internal/models"
)

func setupTestDB(t *testing.T) (*pgx.Conn, func()) {
	t.Helper()

	// Use a temporary database for testing
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/viewdata_test?sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to test DB: %v", err)
	}

	// Clean up tables before test
	_, err = conn.Exec(context.Background(), "DROP TABLE IF EXISTS transactions CASCADE")
	require.NoError(t, err)
	_, err = conn.Exec(context.Background(), "DROP TABLE IF EXISTS statuses CASCADE")
	require.NoError(t, err)

	cleanup := func() {
		conn.Close(context.Background())
	}

	return conn, cleanup
}

func seedTestDB(t *testing.T, conn *pgx.Conn) {
	t.Helper()

	// Insert statuses
	_, err := conn.Exec(context.Background(),
		"INSERT INTO statuses (id, name) VALUES ($1, $2), ($3, $4)",
		0, "SUCCESS", 1, "FAILED")
	require.NoError(t, err)

	// Insert one transaction
	_, err = conn.Exec(context.Background(),
		"INSERT INTO transactions (id, product_id, product_name, amount_text, customer_name, status_id, transaction_date, create_by, create_on) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		int64(1372), "10001", "Test 1", "1000", "abc", 0,
		time.Date(2022, 7, 10, 11, 14, 52, 0, time.UTC),
		"abc", time.Date(2022, 7, 10, 11, 14, 52, 0, time.UTC))
	require.NoError(t, err)
}

func TestMain(m *testing.M) {
	// Set up test DB before running tests
	os.Exit(m.Run())
}

func TestGetViewData_HappyPath(t *testing.T) {
	conn, cleanup := setupTestDB(t)
	defer cleanup()
	seedTestDB(t, conn)

	h := handler.ViewDataHandler{DB: conn}
	req := httptest.NewRequest(http.MethodGet, "/api/view-data", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp models.ViewDataResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)

	assert.Len(t, resp.Data, 1)
	assert.Equal(t, int64(1372), resp.Data[0].ID)
	assert.Equal(t, "10001", resp.Data[0].ProductID)
	assert.Equal(t, "Test 1", resp.Data[0].ProductName)
	assert.Equal(t, "1000", resp.Data[0].Amount)
	assert.Equal(t, "abc", resp.Data[0].CustomerName)
	assert.Equal(t, 0, resp.Data[0].Status)
	assert.Equal(t, "2022-07-10 11:14:52", resp.Data[0].TransactionDate)
	assert.Equal(t, "abc", resp.Data[0].CreateBy)
	assert.Equal(t, "2022-07-10 11:14:52", resp.Data[0].CreateOn)

	assert.Len(t, resp.Status, 2)
	assert.Contains(t, resp.Status, models.StatusResponse{ID: 0, Name: "SUCCESS"})
	assert.Contains(t, resp.Status, models.StatusResponse{ID: 1, Name: "FAILED"})
}

func TestGetViewData_EmptyDB_SeedsData(t *testing.T) {
	conn, cleanup := setupTestDB(t)
	defer cleanup()

	// Ensure tables exist but are empty
	_, err := conn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS statuses (id INTEGER PRIMARY KEY, name TEXT NOT NULL)")
	require.NoError(t, err)
	_, err = conn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS transactions (id BIGINT PRIMARY KEY, product_id TEXT NOT NULL, product_name TEXT NOT NULL, amount_text TEXT NOT NULL, customer_name TEXT NOT NULL, status_id INTEGER NOT NULL REFERENCES statuses(id), transaction_date TIMESTAMP NOT NULL, create_by TEXT NOT NULL, create_on TIMESTAMP NOT NULL)")
	require.NoError(t, err)

	// Verify empty
	count, err := db.CountTransactions(context.Background(), conn)
	require.NoError(t, err)
	assert.Equal(t, 0, count)

	// Manually seed (simulating startup behavior)
	err = db.SeedIfEmpty(context.Background(), conn)
	require.NoError(t, err)

	// Verify seeded
	count, err = db.CountTransactions(context.Background(), conn)
	require.NoError(t, err)
	assert.Equal(t, 12, count) // 12 rows in viewData.json

	// Now test the handler
	h := handler.ViewDataHandler{DB: conn}
	req := httptest.NewRequest(http.MethodGet, "/api/view-data", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp models.ViewDataResponse
	err = json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)

	assert.Len(t, resp.Data, 12)
	assert.Len(t, resp.Status, 2)
}

func TestGetViewData_PartiallySeeded_DB(t *testing.T) {
	conn, cleanup := setupTestDB(t)
	defer cleanup()

	// Pre-populate statuses only
	_, err := conn.Exec(context.Background(),
		"INSERT INTO statuses (id, name) VALUES ($1, $2), ($3, $4)",
		0, "SUCCESS", 1, "FAILED")
	require.NoError(t, err)

	// Ensure transactions table exists but is empty
	_, err = conn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS transactions (id BIGINT PRIMARY KEY, product_id TEXT NOT NULL, product_name TEXT NOT NULL, amount_text TEXT NOT NULL, customer_name TEXT NOT NULL, status_id INTEGER NOT NULL REFERENCES statuses(id), transaction_date TIMESTAMP NOT NULL, create_by TEXT NOT NULL, create_on TIMESTAMP NOT NULL)")
	require.NoError(t, err)

	// Seed only transactions (status already exists)
	err = db.SeedIfEmpty(context.Background(), conn)
	require.NoError(t, err)

	// Verify transactions were seeded
	count, err := db.CountTransactions(context.Background(), conn)
	require.NoError(t, err)
	assert.Equal(t, 12, count)

	// Handler should return both data and status
	h := handler.ViewDataHandler{DB: conn}
	req := httptest.NewRequest(http.MethodGet, "/api/view-data", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp models.ViewDataResponse
	err = json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)

	assert.Len(t, resp.Data, 12)
	assert.Len(t, resp.Status, 2)
}

func TestGetViewData_InvalidDatetimeFormat(t *testing.T) {
	// This test verifies that malformed datetime in JSON causes startup failure.
	// Since seeding happens at app startup, we simulate by calling seed with bad data.
	conn, cleanup := setupTestDB(t)
	defer cleanup()

	// Create a temporary JSON file with malformed datetime
	tmpFile, err := os.CreateTemp("", "bad-datetime-*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	badJSON := `{
		"data": [
			{
				"id": 1372,
				"productID": "10001",
				"productName": "Test 1",
				"amount": "1000",
				"customerName": "abc",
				"status": 0,
				"transactionDate": "not-a-date",
				"createBy": "abc",
				"createOn": "2022-07-10 11:14:52"
			}
		],
		"status": [
			{"id": 0, "name": "SUCCESS"},
			{"id": 1, "name": "FAILED"}
		]
	}`
	_, err = tmpFile.Write([]byte(badJSON))
	require.NoError(t, err)
	tmpFile.Close()

	// Temporarily override the embedded JSON path
	originalJSONPath := db.EmbeddedJSONPath
	db.EmbeddedJSONPath = tmpFile.Name()
	defer func() { db.EmbeddedJSONPath = originalJSONPath }()

	// Insert statuses first
	_, err = conn.Exec(context.Background(),
		"INSERT INTO statuses (id, name) VALUES ($1, $2), ($3, $4)",
		0, "SUCCESS", 1, "FAILED")
	require.NoError(t, err)

	// Seeding should fail due to malformed datetime
	err = db.SeedIfEmpty(context.Background(), conn)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "parsing time")
}

func TestGetViewData_MissingStatusReference(t *testing.T) {
	conn, cleanup := setupTestDB(t)
	defer cleanup()

	// Create a temporary JSON file with invalid status reference
	tmpFile, err := os.CreateTemp("", "bad-status-*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	badJSON := `{
		"data": [
			{
				"id": 1372,
				"productID": "10001",
				"productName": "Test 1",
				"amount": "1000",
				"customerName": "abc",
				"status": 999,
				"transactionDate": "2022-07-10 11:14:52",
				"createBy": "abc",
				"createOn": "2022-07-10 11:14:52"
			}
		],
		"status": [
			{"id": 0, "name": "SUCCESS"},
			{"id": 1, "name": "FAILED"}
		]
	}`
	_, err = tmpFile.Write([]byte(badJSON))
	require.NoError(t, err)
	tmpFile.Close()

	// Temporarily override the embedded JSON path
	originalJSONPath := db.EmbeddedJSONPath
	db.EmbeddedJSONPath = tmpFile.Name()
	defer func() { db.EmbeddedJSONPath = originalJSONPath }()

	// Seeding should skip invalid transaction
	err = db.SeedIfEmpty(context.Background(), conn)
	require.NoError(t, err)

	// Verify no transactions were inserted (invalid ones skipped)
	count, err := db.CountTransactions(context.Background(), conn)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestGetViewData_DBConnectionFailure(t *testing.T) {
	// Simulate DB connection failure by using invalid credentials
	cfg := db.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "invalid_user",
		Password: "wrong_password",
		Database: "viewdata",
	}

	_, err := db.Connect(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "connection refused") || assert.Contains(t, err.Error(), "password authentication failed")
}

func TestGetViewData_DuplicatePrimaryKeys(t *testing.T) {
	conn, cleanup := setupTestDB(t)
	defer cleanup()

	// Create a temporary JSON file with duplicate IDs (last one wins)
	tmpFile, err := os.CreateTemp("", "dup-id-*.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	dupJSON := `{
		"data": [
			{
				"id": 1372,
				"productID": "10001",
				"productName": "Test 1",
				"amount": "1000",
				"customerName": "abc",
				"status": 0,
				"transactionDate": "2022-07-10 11:14:52",
				"createBy": "abc",
				"createOn": "2022-07-10 11:14:52"
			},
			{
				"id": 1372,
				"productID": "10002",
				"productName": "Test 2",
				"amount": "2000",
				"customerName": "xyz",
				"status": 1,
				"transactionDate": "2022-07-11 13:14:52",
				"createBy": "xyz",
				"createOn": "2022-07-11 13:14:52"
			}
		],
		"status": [
			{"id": 0, "name": "SUCCESS"},
			{"id": 1, "name": "FAILED"}
		]
	}`
	_, err = tmpFile.Write([]byte(dupJSON))
	require.NoError(t, err)
	tmpFile.Close()

	// Temporarily override the embedded JSON path
	originalJSONPath := db.EmbeddedJSONPath
	db.EmbeddedJSONPath = tmpFile.Name()
	defer func() { db.EmbeddedJSONPath = originalJSONPath }()

	// Insert statuses first
	_, err = conn.Exec(context.Background(),
		"INSERT INTO statuses (id, name) VALUES ($1, $2), ($3, $4)",
		0, "SUCCESS", 1, "FAILED")
	require.NoError(t, err)

	// Seed should succeed, last occurrence wins
	err = db.SeedIfEmpty(context.Background(), conn)
	require.NoError(t, err)

	// Verify only one transaction with ID 1372, with last values
	var tx models.Transaction
	err = conn.QueryRow(context.Background(),
		"SELECT id, product_id, product_name, amount_text, customer_name, status_id, transaction_date, create_by, create_on FROM transactions WHERE id = $1", 1372).
		Scan(&tx.ID, &tx.ProductID, &tx.ProductName, &tx.AmountText, &tx.CustomerName, &tx.StatusID, &tx.TransactionDate, &tx.CreateBy, &tx.CreateOn)
	require.NoError(t, err)

	assert.Equal(t, int64(1372), tx.ID)
	assert.Equal(t, "10002", tx.ProductID) // last occurrence
	assert.Equal(t, "Test 2", tx.ProductName)
	assert.Equal(t, "2000", tx.AmountText)
	assert.Equal(t, "xyz", tx.CustomerName)
	assert.Equal(t, 1, tx.StatusID)
}

func TestGetViewData_LargeResponsePayload(t *testing.T) {
	conn, cleanup := setupTestDB(t)
	defer cleanup()

	// Insert 10,000+ transactions
	_, err := conn.Exec(context.Background(),
		"INSERT INTO statuses (id, name) VALUES ($1, $2), ($3, $4)",
		0, "SUCCESS", 1, "FAILED")
	require.NoError(t, err)

	batchSize := 1000
	totalRows := 11000
	for i := 0; i < totalRows; i += batchSize {
		var args []interface{}
		vals := make([]string, 0, batchSize*9)
		for j := 0; j < batchSize; j++ {
			id := int64(2000 + i + j)
			vals = append(vals, "$"+string(rune(49+len(vals))), "$"+string(rune(49+len(vals))), "$"+string(rune(49+len(vals))), "$"+string(rune(49+len(vals))), "$"+string(rune(49+len(vals))), "$"+string(rune(49+len(vals))), "$"+string(rune(49+len(vals))), "$"+string(rune(49+len(vals))), "$"+string(rune(49+len(vals)))))
			args = append(args, id, "PROD", "Product", "100", "user", 0,
				time.Now(), "admin", time.Now())
		}
		_, err = conn.Exec(context.Background(),
			"INSERT INTO transactions (id, product_id, product_name, amount_text, customer_name, status_id, transaction_date, create_by, create_on) VALUES "+vals[0]+", "+vals[1]+", "+vals[2]+", "+vals[3]+", "+vals[4]+", "+vals[5]+", "+vals[6]+", "+vals[7]+", "+vals[8], args...)
		require.NoError(t, err)
	}

	// Verify count
	count, err := db.CountTransactions(context.Background(), conn)
	require.NoError(t, err)
	assert.Equal(t, totalRows, count)

	// Test handler can stream large response
	h := handler.ViewDataHandler{DB: conn}
	req := httptest.NewRequest(http.MethodGet, "/api/view-data", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Count rows in response
	var resp models.ViewDataResponse
	err = json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Len(t, resp.Data, totalRows)
}
