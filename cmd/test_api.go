package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"internship-test-api/internal/model"
	"internship-test-api/internal/repository"
	"internship-test-api/internal/service"
)

func TestAPI(t *testing.T) {
	// Setup in-memory DB for testing
	db, err := setupTestDB()
	if err != nil {
		t.Fatal(err)
	}

	statusRepo := repository.NewStatusRepo(db)
	transactionRepo := repository.NewTransactionRepo(db)
	viewDataService := service.NewViewDataService(transactionRepo, statusRepo)

	router := setupRouter(viewDataService)

	t.Run("GET /api/view-data returns 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/view-data", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}

		var resp ViewDataResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if len(resp.Data) == 0 {
			t.Error("expected non-empty data array")
		}
		if len(resp.Status) == 0 {
			t.Error("expected non-empty status array")
		}
	})
}

func setupRouter(viewDataService service.ViewDataService) *http.Router {
	router := http.NewRouter()
	router.HandleFunc("/api/view-data", handler.GetViewData(viewDataService)).Methods(http.MethodGet)
	return router
}

func setupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open("host=localhost user=test password=test dbname=test_viewdata port=5432 sslmode=disable TimeZone=UTC"))
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&model.Transaction{}, &model.Status{}); err != nil {
		return nil, err
	}
	return db, nil
}

type ViewDataResponse struct {
	Data   []ViewDataResponseRow `json:"data"`
	Status []StatusItem          `json:"status"`
}

type ViewDataResponseRow struct {
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

type StatusItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
