package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTransactionModel(t *testing.T) {
	transactionDate := time.Date(2023, 1, 15, 9, 30, 0, 0, time.UTC)
	createOn := time.Date(2023, 1, 15, 9, 30, 0, 0, time.UTC)
	
	transaction := Transaction{
		ID:              1001,
		ProductID:       "P001",
		ProductName:     "Test Product",
		AmountText:      "1234.56",
		CustomerName:    "Test Customer",
		StatusID:        0,
		TransactionDate: transactionDate,
		CreateBy:        "testuser",
		CreateOn:        createOn,
	}
	
	assert.Equal(t, int64(1001), transaction.ID)
	assert.Equal(t, "P001", transaction.ProductID)
	assert.Equal(t, "Test Product", transaction.ProductName)
	assert.Equal(t, "1234.56", transaction.AmountText)
	assert.Equal(t, "Test Customer", transaction.CustomerName)
	assert.Equal(t, 0, transaction.StatusID)
	assert.Equal(t, transactionDate, transaction.TransactionDate)
	assert.Equal(t, "testuser", transaction.CreateBy)
	assert.Equal(t, createOn, transaction.CreateOn)
}

func TestStatusModel(t *testing.T) {
	status := Status{
		ID:   1,
		Name: "TEST_STATUS",
	}
	
	assert.Equal(t, 1, status.ID)
	assert.Equal(t, "TEST_STATUS", status.Name)
}

func TestTransactionRowConversion(t *testing.T) {
	transactionDate := time.Date(2023, 1, 15, 9, 30, 0, 0, time.UTC)
	createOn := time.Date(2023, 1, 15, 9, 30, 0, 0, time.UTC)
	
	transaction := Transaction{
		ID:              1001,
		ProductID:       "P001",
		ProductName:     "Test Product",
		AmountText:      "1234.56",
		CustomerName:    "Test Customer",
		StatusID:        1,
		TransactionDate: transactionDate,
		CreateBy:        "testuser",
		CreateOn:        createOn,
	}
	
	row := TransactionRow{
		ID:              transaction.ID,
		ProductID:       transaction.ProductID,
		ProductName:     transaction.ProductName,
		Amount:          transaction.AmountText,
		CustomerName:    transaction.CustomerName,
		Status:          transaction.StatusID,
		TransactionDate: transaction.TransactionDate.Format("2006-01-02 15:04:05"),
		CreateBy:        transaction.CreateBy,
		CreateOn:        transaction.CreateOn.Format("2006-01-02 15:04:05"),
	}
	
	assert.Equal(t, int64(1001), row.ID)
	assert.Equal(t, "P001", row.ProductID)
	assert.Equal(t, "Test Product", row.ProductName)
	assert.Equal(t, "1234.56", row.Amount)
	assert.Equal(t, "Test Customer", row.CustomerName)
	assert.Equal(t, 1, row.Status)
	assert.Equal(t, "2023-01-15 09:30:00", row.TransactionDate)
	assert.Equal(t, "testuser", row.CreateBy)
	assert.Equal(t, "2023-01-15 09:30:00", row.CreateOn)
}

func TestStatusItemConversion(t *testing.T) {
	status := Status{
		ID:   2,
		Name: "CONVERTED",
	}
	
	item := StatusItem{
		ID:   status.ID,
		Name: status.Name,
	}
	
	assert.Equal(t, 2, item.ID)
	assert.Equal(t, "CONVERTED", item.Name)
}
