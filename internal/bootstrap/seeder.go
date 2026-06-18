package bootstrap

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/PaulusBilly/internship-test-API/internal/models"
	"github.com/PaulusBilly/internship-test-API/internal/repository"
)

const dateLayout = "2006-01-02 15:04:05"

type ViewDataFile struct {
	Data   []Row   `json:"data"`
	Status []StItem `json:"status"`
}

type Row struct {
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

type StItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func Seed(repo repository.Repository, jsonPath string) error {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to read viewData.json: %w", err)
	}

	var file ViewDataFile
	if err := json.Unmarshal(data, &file); err != nil {
		return fmt.Errorf("failed to parse viewData.json: %w", err)
	}

	// Seed statuses
	for _, st := range file.Status {
		if err := repo.SaveStatus(models.Status{ID: st.ID, Name: st.Name}); err != nil {
			return fmt.Errorf("failed to save status %d: %w", st.ID, err)
		}
	}

	// Seed transactions
	for _, r := range file.Data {
		// Validate status exists
		found := false
		for _, st := range file.Status {
			if st.ID == r.Status {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("missing status id %d referenced in transaction %d", r.Status, r.ID)
		}

		txnDate, err := time.Parse(dateLayout, r.TransactionDate)
		if err != nil {
			return fmt.Errorf("invalid transactionDate format for id %d: %w", r.ID, err)
		}
		createOn, err := time.Parse(dateLayout, r.CreateOn)
		if err != nil {
			return fmt.Errorf("invalid createOn format for id %d: %w", r.ID, err)
		}

		t := models.Transaction{
			ID:              r.ID,
			ProductID:       r.ProductID,
			ProductName:     r.ProductName,
			AmountText:      r.Amount,
			CustomerName:    r.CustomerName,
			StatusID:        r.Status,
			TransactionDate: txnDate,
			CreateBy:        r.CreateBy,
			CreateOn:        createOn,
		}
		if err := repo.SaveTransaction(t); err != nil {
			return fmt.Errorf("failed to save transaction %d: %w", r.ID, err)
		}
	}

	return nil
}
