package service

import (
	"testing"
	"time"

	"github.com/PaulusBilly/internship-test-API/internal/dto"
	"github.com/PaulusBilly/internship-test-API/internal/model"
	"github.com/PaulusBilly/internship-test-API/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) FindAll() ([]model.Transaction, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Transaction), args.Error(1)
}

type MockStatusRepository struct {
	mock.Mock
}

func (m *MockStatusRepository) FindAll() ([]model.Status, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Status), args.Error(1)
}

func TestViewDataService_BuildResponse(t *testing.T) {
	t.Run("successful build with data", func(t *testing.T) {
		mockTxRepo := new(MockTransactionRepository)
		mockStatusRepo := new(MockStatusRepository)

		status0 := model.Status{ID: 0, Name: "SUCCESS"}
		status1 := model.Status{ID: 1, Name: "FAILED"}
		transactionTime, _ := time.Parse("2006-01-02 15:04:05", "2022-07-10 11:14:52")

		transactions := []model.Transaction{
			{
				ID:              1372,
				ProductID:       "10001",
				ProductName:     "Test 1",
				AmountText:      "1000",
				CustomerName:    "abc",
				Status:          status0,
				TransactionDate: transactionTime,
				CreateBy:        "abc",
				CreateOn:        transactionTime,
			},
		}
		statuses := []model.Status{status0, status1}

		mockTxRepo.On("FindAll").Return(transactions, nil)
		mockStatusRepo.On("FindAll").Return(statuses, nil)

		svc := NewViewDataService(mockTxRepo, mockStatusRepo)
		resp, err := svc.BuildResponse()

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Data, 1)
		assert.Len(t, resp.Status, 2)

		row := resp.Data[0]
		assert.Equal(t, int64(1372), row.ID)
		assert.Equal(t, "10001", row.ProductID)
		assert.Equal(t, "Test 1", row.ProductName)
		assert.Equal(t, "1000", row.Amount)
		assert.Equal(t, "abc", row.CustomerName)
		assert.Equal(t, 0, row.Status)
		assert.Equal(t, "2022-07-10 11:14:52", row.TransactionDate)
		assert.Equal(t, "abc", row.CreateBy)
		assert.Equal(t, "2022-07-10 11:14:52", row.CreateOn)

		statusItem0 := resp.Status[0]
		assert.Equal(t, 0, statusItem0.ID)
		assert.Equal(t, "SUCCESS", statusItem0.Name)

		statusItem1 := resp.Status[1]
		assert.Equal(t, 1, statusItem1.ID)
		assert.Equal(t, "FAILED", statusItem1.Name)

		mockTxRepo.AssertExpectations(t)
		mockStatusRepo.AssertExpectations(t)
	})

	t.Run("transaction repository error", func(t *testing.T) {
		mockTxRepo := new(MockTransactionRepository)
		mockStatusRepo := new(MockStatusRepository)

		mockTxRepo.On("FindAll").Return(nil, assert.AnError)

		svc := NewViewDataService(mockTxRepo, mockStatusRepo)
		resp, err := svc.BuildResponse()

		assert.Error(t, err)
		assert.Nil(t, resp)
		mockTxRepo.AssertExpectations(t)
		mockStatusRepo.AssertNotCalled(t, "FindAll")
	})

	t.Run("status repository error", func(t *testing.T) {
		mockTxRepo := new(MockTransactionRepository)
		mockStatusRepo := new(MockStatusRepository)

		transactions := []model.Transaction{}
		mockTxRepo.On("FindAll").Return(transactions, nil)
		mockStatusRepo.On("FindAll").Return(nil, assert.AnError)

		svc := NewViewDataService(mockTxRepo, mockStatusRepo)
		resp, err := svc.BuildResponse()

		assert.Error(t, err)
		assert.Nil(t, resp)
		mockTxRepo.AssertExpectations(t)
		mockStatusRepo.AssertExpectations(t)
	})

	t.Run("empty data", func(t *testing.T) {
		mockTxRepo := new(MockTransactionRepository)
		mockStatusRepo := new(MockStatusRepository)

		transactions := []model.Transaction{}
		statuses := []model.Status{}

		mockTxRepo.On("FindAll").Return(transactions, nil)
		mockStatusRepo.On("FindAll").Return(statuses, nil)

		svc := NewViewDataService(mockTxRepo, mockStatusRepo)
		resp, err := svc.BuildResponse()

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Empty(t, resp.Data)
		assert.Empty(t, resp.Status)

		mockTxRepo.AssertExpectations(t)
		mockStatusRepo.AssertExpectations(t)
	})
}
