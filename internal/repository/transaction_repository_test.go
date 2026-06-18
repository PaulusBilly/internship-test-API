package repository

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"internship-test-api/internal/model"
)

func TestTransactionRepository_FindAll(t *testing.T) {
	db, err := sql.Open("postgres", os.Getenv("TEST_DB_URL"))
	if err != nil {
		t.Fatalf("failed to open DB: %v", err)
	}
	defer db.Close()

	repo := NewTransactionRepository(db)

	txs, err := repo.FindAll()
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}

	if len(txs) == 0 {
		t.Errorf("expected at least one transaction, got %d", len(txs))
	}

	for _, tx := range txs {
		if tx.ID == 0 {
			t.Error("transaction ID is zero")
		}
		if tx.ProductID == "" {
			t.Error("transaction ProductID is empty")
		}
	}
}
