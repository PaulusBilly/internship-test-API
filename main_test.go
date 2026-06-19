package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestViewDataHandler_Success(t *testing.T) {
	// Setup: load and embed the real viewData.json
	jsonPath := filepath.Join("data", "viewData.json")
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		t.Fatalf("failed to read viewData.json: %v", err)
	}

	var expected ViewDataResponse
	if err := json.Unmarshal(data, &expected); err != nil {
		t.Fatalf("failed to unmarshal viewData.json: %v", err)
	}

	// Act
	req := httptest.NewRequest(http.MethodGet, "/api/view-data", nil)
	w := httptest.NewRecorder()
	viewDataHandler(w, req)
	resp := w.Result()

	// Assert status and content-type
	if got := resp.StatusCode; got != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, got)
	}
	if got := resp.Header.Get("Content-Type"); got != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got %q", got)
	}

	// Assert body matches expected JSON structure
	var actual ViewDataResponse
	if err := json.NewDecoder(resp.Body).Decode(&actual); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	// Compare data length and first item for structure sanity
	if len(actual.Data) != len(expected.Data) {
		t.Errorf("data length mismatch: expected %d, got %d", len(expected.Data), len(actual.Data))
	}
	if len(actual.Status) != len(expected.Status) {
		t.Errorf("status length mismatch: expected %d, got %d", len(expected.Status), len(actual.Status))
	}

	// Validate first data item fields
	if len(actual.Data) > 0 && len(expected.Data) > 0 {
		got := actual.Data[0]
		want := expected.Data[0]
		if got.ID != want.ID {
			t.Errorf("first data.ID: expected %d, got %d", want.ID, got.ID)
		}
		if got.ProductID != want.ProductID {
			t.Errorf("first data.productID: expected %q, got %q", want.ProductID, got.ProductID)
		}
		if got.ProductName != want.ProductName {
			t.Errorf("first data.productName: expected %q, got %q", want.ProductName, got.ProductName)
		}
		if got.Amount != want.Amount {
			t.Errorf("first data.amount: expected %q, got %q", want.Amount, got.Amount)
		}
		if got.CustomerName != want.CustomerName {
			t.Errorf("first data.customerName: expected %q, got %q", want.CustomerName, got.CustomerName)
		}
		if got.Status != want.Status {
			t.Errorf("first data.status: expected %d, got %d", want.Status, got.Status)
		}
		if got.TransactionDate != want.TransactionDate {
			t.Errorf("first data.transactionDate: expected %q, got %q", want.TransactionDate, got.TransactionDate)
		}
		if got.CreateBy != want.CreateBy {
			t.Errorf("first data.createBy: expected %q, got %q", want.CreateBy, got.CreateBy)
		}
		if got.CreateOn != want.CreateOn {
			t.Errorf("first data.createOn: expected %q, got %q", want.CreateOn, got.CreateOn)
		}
	}

	// Validate status items
	if len(actual.Status) > 0 && len(expected.Status) > 0 {
		got := actual.Status[0]
		want := expected.Status[0]
		if got.ID != want.ID {
			t.Errorf("first status.id: expected %d, got %d", want.ID, got.ID)
		}
		if got.Name != want.Name {
			t.Errorf("first status.name: expected %q, got %q", want.Name, got.Name)
		}
	}
}

func TestViewDataHandler_EdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		setup      func()
		wantStatus int
	}{
		{
			name: "E3 - GET /api/view-data returns 200",
			setup: func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			req := httptest.NewRequest(http.MethodGet, "/api/view-data", nil)
			w := httptest.NewRecorder()
			viewDataHandler(w, req)
			resp := w.Result()

			if got := resp.StatusCode; got != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, got)
			}
		})
	}
}
