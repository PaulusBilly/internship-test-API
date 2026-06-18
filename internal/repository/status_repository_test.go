package repository

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"internship-test-api/internal/model"
)

func TestStatusRepository_FindAll(t *testing.T) {
	db, err := sql.Open("postgres", os.Getenv("TEST_DB_URL"))
	if err != nil {
		t.Fatalf("failed to open DB: %v", err)
	}
	defer db.Close()

	repo := NewStatusRepository(db)

	statuses, err := repo.FindAll()
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}

	if len(statuses) == 0 {
		t.Errorf("expected at least one status, got %d", len(statuses))
	}

	for _, st := range statuses {
		if st.ID == 0 {
			t.Error("status ID is zero")
		}
		if st.Name == "" {
			t.Error("status Name is empty")
		}
	}
}
