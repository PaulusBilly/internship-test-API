package service

import (
	"testing"
	"time"

	"internship-test-API/internal/dto"
	"internship-test-API/internal/models"
	"internship-test-API/internal/repository"
)

type MockRepo struct {
	transactions []models.Transaction
	statuses     []models.Status
	err          error
}

func (m *MockRepo) GetTransactions() ([]models.Transaction, error) {
	return m.transactions, m.err
}

func (m *MockRepo) GetStatuses() ([]models.Status, error) {
	return m.statuses, m.err
}

func (m *MockRepo) SaveStatus(s models.Status) error {
	return m.err
}

func (m *MockRepo) SaveTransaction(t models.Transaction) error {
	return m.err
}

func (m *MockRepo) CountTransactions() (int64, error) {
	return int64(len(m.transactions)), m.err
}

func TestBuildResponse_HappyPath(t *testing.T) {
	mockRepo := &MockRepo{
		transactions: []models.Transaction{
			{
				ID:              1372,
				ProductID:       "10001",
				ProductName:     "Test 1",
				AmountText:      "1000",
				CustomerName:    "abc",
				StatusID:        0,
				TransactionDate: time.Date(2022, 7, 10, 11, 14, 52, 0, time.UTC),
				CreateBy:        "abc",
				CreateOn:        time.Date(2022, 7, 10, 11, 14, 52, 0, time.UTC),
			},
		},
		statuses: []models.Status{
			{ID: 0, Name: "SUCCESS"},
			{ID: 1, Name: "FAILED"},
		},
	}

	service := &Service{repo: mockRepo}
	resp, err := service.BuildResponse()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(resp.Data) != 1 {
		t.Errorf("Expected 1 data row, got %d", len(resp.Data))
	}

	if len(resp.Status) != 2 {
		t.Errorf("Expected 2 status items, got %d", len(resp.Status))
	}

	row := resp.Data[0]
	if row.ProductID != "10001" {
		t.Errorf("Expected ProductID '10001', got %s", row.ProductID)
	}
	if row.Amount != "1000" {
		t.Errorf("Expected Amount '1000', got %s", row.Amount)
	}
	if row.Status != 0 {
		t.Errorf("Expected Status 0, got %d", row.Status)
	}
	if row.TransactionDate != "2022-07-10 11:14:52" {
		t.Errorf("Expected TransactionDate '2022-07-10 11:14:52', got %s", row.TransactionDate)
	}
}

func TestBuildResponse_EmptyData(t *testing.T) {
	mockRepo := &MockRepo{
		transactions: []models.Transaction{},
		statuses:     []models.Status{},
	}

	service := &Service{repo: mockRepo}
	resp, err := service.BuildResponse()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(resp.Data) != 0 {
		t.Errorf("Expected 0 data rows, got %d", len(resp.Data))
	}
	if len(resp.Status) != 0 {
		t.Errorf("Expected 0 status items, got %d", len(resp.Status))
	}
}

func TestBuildResponse_RepoError(t *testing.T) {
	mockRepo := &MockRepo{
		err: repository.ErrDatabase,
	}

	service := &Service{repo: mockRepo}
	_, err := service.BuildResponse()
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}
