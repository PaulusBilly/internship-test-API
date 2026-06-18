package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"internship-test-api/internal/model"
)

type SeedService struct {
	db *sql.DB
}

func NewSeedService(db *sql.DB) *SeedService {
	return &SeedService{db: db}
}

func (s *SeedService) Seed() error {
	// Check if data already exists
	var count int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM transactions`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil // already seeded
	}

	// Open JSON file
	f, err := os.Open("src/main/resources/viewData.json")
	if err != nil {
		return fmt.Errorf("failed to open viewData.json: %w", err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("failed to read viewData.json: %w", err)
	}

	var file struct {
		Data   []jsonRow   `json:"data"`
		Status []jsonStatus `json:"status"`
	}

	if err := json.Unmarshal(data, &file); err != nil {
		return fmt.Errorf("failed to parse viewData.json: %w", err)
	}

	// Seed statuses
	for _, st := range file.Status {
		if _, err := s.db.ExecContext(context.Background(),
			`INSERT INTO statuses (id, name) VALUES ($1, $2) ON CONFLICT (id) DO NOTHING`,
			st.ID, st.Name); err != nil {
			return fmt.Errorf("failed to seed status %d: %w", st.ID, err)
		}
	}

	// Seed transactions
	for _, r := range file.Data {
		txnDate, err := time.Parse("2006-01-02 15:04:05", r.TransactionDate)
		if err != nil {
			return fmt.Errorf("invalid transaction date %s: %w", r.TransactionDate, err)
		}
		createOn, err := time.Parse("2006-01-02 15:04:05", r.CreateOn)
		if err != nil {
			return fmt.Errorf("invalid createOn date %s: %w", r.CreateOn, err)
		}

		if _, err := s.db.ExecContext(context.Background(),
			`INSERT INTO transactions (
				id, product_id, product_name, amount_text, customer_name,
				status_id, transaction_date, create_by, create_on
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			ON CONFLICT (id) DO NOTHING`,
			r.ID, r.ProductID, r.ProductName, r.Amount, r.CustomerName,
			r.Status, txnDate, r.CreateBy, createOn); err != nil {
			return fmt.Errorf("failed to seed transaction %d: %w", r.ID, err)
		}
	}

	return nil
}

type jsonRow struct {
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

type jsonStatus struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
