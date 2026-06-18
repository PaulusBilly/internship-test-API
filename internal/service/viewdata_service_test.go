package service

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"internship-test-api/internal/handler"
	"internship-test-api/internal/repository"
)

func TestViewDataService_GetViewData(t *testing.T) {
	db, err := sql.Open("postgres", os.Getenv("TEST_DB_URL"))
	if err != nil {
		t.Fatalf("failed to open DB: %v", err)
	}
	defer db.Close()

	statusRepo := repository.NewStatusRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	svc := NewViewDataService(statusRepo, transactionRepo)

	resp, err := svc.GetViewData()
	if err != nil {
		t.Fatalf("GetViewData failed: %v", err)
	}

	if len(resp.Data) == 0 {
		t.Errorf("expected at least one transaction, got %d", len(resp.Data))
	}

	if len(resp.Status) == 0 {
		t.Errorf("expected at least one status, got %d", len(resp.Status))
	}

	for _, row := range resp.Data {
		if row.ID == 0 {
			t.Error("row ID is zero")
		}
		if row.ProductID == "" {
			t.Error("row ProductID is empty")
		}
		if row.Status < 0 {
			t.Error("row Status is negative")
		}
	}
}
