package bootstrap

import (
	"errors"
	"testing"

	"internship-test-API/internal/models"
	"internship-test-API/internal/repository"
)

type MockRepo struct {
	count          int64
	saveStatusErr  error
	saveTxErr      error
	countErr       error
	savedStatuses  []models.Status
	savedTxs       []models.Transaction
}

func (m *MockRepo) GetTransactions() ([]models.Transaction, error) {
	return nil, nil
}

func (m *MockRepo) GetStatuses() ([]models.Status, error) {
	return nil, nil
}

func (m *MockRepo) SaveStatus(s models.Status) error {
	m.savedStatuses = append(m.savedStatuses, s)
	return m.saveStatusErr
}

func (m *MockRepo) SaveTransaction(t models.Transaction) error {
	m.savedTxs = append(m.savedTxs, t)
	return m.saveTxErr
}

func (m *MockRepo) CountTransactions() (int64, error) {
	return m.count, m.countErr
}

func TestSeed_E1_EmptyDatabase(t *testing.T) {
	mockRepo := &MockRepo{count: 0}
	err := Seed(mockRepo, "../../resources/viewData.json")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(mockRepo.savedStatuses) != 2 {
		t.Errorf("Expected 2 statuses saved, got %d", len(mockRepo.savedStatuses))
	}
	if len(mockRepo.savedTxs) == 0 {
		t.Error("Expected transactions to be saved, got 0")
	}
}

func TestSeed_E2_ExistingData(t *testing.T) {
	mockRepo := &MockRepo{count: 1}
	err := Seed(mockRepo, "../../resources/viewData.json")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(mockRepo.savedStatuses) > 0 {
		t.Error("Expected no statuses saved when data exists")
	}
	if len(mockRepo.savedTxs) > 0 {
		t.Error("Expected no transactions saved when data exists")
	}
}

func TestSeed_E3_MissingStatusReference(t *testing.T) {
	// Create test JSON with invalid status reference
	testJSON := `{
		"data": [{
			"id": 999,
			"productID": "test",
			"productName": "test",
			"amount": "100",
			"customerName": "test",
			"status": 99,
			"transactionDate": "2022-01-01 00:00:00",
			"createBy": "test",
			"createOn": "2022-01-01 00:00:00"
		}],
		"status": [{"id": 0, "name": "SUCCESS"}]
	}`

	// Write to temporary file and test
	tmpFile := createTempJSON(t, testJSON)
	defer tmpFile.Close()

	mockRepo := &MockRepo{count: 0}
	err := Seed(mockRepo, tmpFile.Name())
	if err == nil {
		t.Fatal("Expected error for missing status reference, got nil")
	}
}

func TestSeed_E4_DateFormatMismatch(t *testing.T) {
	testJSON := `{
		"data": [{
			"id": 999,
			"productID": "test",
			"productName": "test",
			"amount": "100",
			"customerName": "test",
			"status": 0,
			"transactionDate": "invalid-date",
			"createBy": "test",
			"createOn": "2022-01-01 00:00:00"
		}],
		"status": [{"id": 0, "name": "SUCCESS"}]
	}`

	tmpFile := createTempJSON(t, testJSON)
	defer tmpFile.Close()

	mockRepo := &MockRepo{count: 0}
	err := Seed(mockRepo, tmpFile.Name())
	if err == nil {
		t.Fatal("Expected error for invalid date format, got nil")
	}
}

func TestSeed_E5_DatabaseError(t *testing.T) {
	mockRepo := &MockRepo{
		count:   0,
		countErr: errors.New("database connection failed"),
	}

	err := Seed(mockRepo, "../../resources/viewData.json")
	if err == nil {
		t.Fatal("Expected error for database failure, got nil")
	}
}
