package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// E2: Happy path — endpoint returns correct JSON structure and status
func TestViewDataHandler_HappyPath(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/view-data", nil)
	rec := httptest.NewRecorder()

	viewDataHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if got := resp.StatusCode; got != http.StatusOK {
		t.Fatalf("expected status 200, got %d", got)
	}

	var actual ViewDataResponse
	if err := json.NewDecoder(resp.Body).Decode(&actual); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(actual.Data) != 12 {
		t.Fatalf("expected 12 data items, got %d", len(actual.Data))
	}
	if len(actual.Status) != 2 {
		t.Fatalf("expected 2 status items, got %d", len(actual.Status))
	}

	// Verify first data item matches expected values
	first := actual.Data[0]
	if first.ID != 1372 {
		t.Errorf("expected first ID 1372, got %d", first.ID)
	}
	if first.ProductID != "10001" {
		t.Errorf("expected first productID '10001', got %q", first.ProductID)
	}
	if first.Status != 0 {
		t.Errorf("expected first status 0, got %d", first.Status)
	}

	// Verify status items
	if actual.Status[0].ID != 0 || actual.Status[0].Name != "SUCCESS" {
		t.Errorf("expected status[0] = {0, SUCCESS}, got %+v", actual.Status[0])
	}
	if actual.Status[1].ID != 1 || actual.Status[1].Name != "FAILED" {
		t.Errorf("expected status[1] = {1, FAILED}, got %+v", actual.Status[1])
	}
}

// E4: Verify Content-Type header
func TestViewDataHandler_ContentType(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/view-data", nil)
	rec := httptest.NewRecorder()

	viewDataHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if got := resp.Header.Get("Content-Type"); !strings.Contains(got, "application/json") {
		t.Fatalf("E4: expected Content-Type 'application/json', got %q", got)
	}
}

// E3: Verify timestamp format
func TestViewDataHandler_TimestampFormat(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/view-data", nil)
	rec := httptest.NewRecorder()

	viewDataHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	var actual ViewDataResponse
	if err := json.NewDecoder(resp.Body).Decode(&actual); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	layout := "2006-01-02 15:04:05"
	for i, row := range actual.Data {
		if _, err := time.Parse(layout, row.TransactionDate); err != nil {
			t.Fatalf("E3: invalid transactionDate at index %d: %q", i, row.TransactionDate)
		}
		if _, err := time.Parse(layout, row.CreateOn); err != nil {
			t.Fatalf("E3: invalid createOn at index %d: %q", i, row.CreateOn)
		}
	}
}

// E5: Verify data completeness
func TestViewDataHandler_DataCompleteness(t *testing.T) {
	if len(viewData.Data) != 12 {
		t.Fatalf("E5: expected 12 data items at startup, got %d", len(viewData.Data))
	}
	if len(viewData.Status) != 2 {
		t.Fatalf("E5: expected 2 status items at startup, got %d", len(viewData.Status))
	}

	req := httptest.NewRequest(http.MethodGet, "/api/view-data", nil)
	rec := httptest.NewRecorder()

	viewDataHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	var actual ViewDataResponse
	if err := json.NewDecoder(resp.Body).Decode(&actual); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(actual.Data) != 12 || len(actual.Status) != 2 {
		t.Fatalf("E5: response mismatch: data=%d, status=%d", len(actual.Data), len(actual.Status))
	}
}

// E1: Startup failure on malformed JSON — tested via manual inspection of main()
// (Cannot be unit-tested directly without refactoring main() — spec requires panic on parse error)

// E6: Port conflict — not testable in unit tests (requires runtime port binding failure simulation)
// Verified via manual inspection of main() error handling

// Additional edge case: Verify exact field names match spec
func TestViewDataHandler_FieldNames(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/view-data", nil)
	rec := httptest.NewRecorder()

	viewDataHandler(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	var raw map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Check root keys
	if _, ok := raw["data"]; !ok {
		t.Fatalf("expected root key 'data'")
	}
	if _, ok := raw["status"]; !ok {
		t.Fatalf("expected root key 'status'")
	}

	// Check first row keys
	dataSlice, ok := raw["data"].([]any)
	if !ok || len(dataSlice) == 0 {
		t.Fatalf("expected non-empty 'data' array")
	}
	firstRow, ok := dataSlice[0].(map[string]any)
	if !ok {
		t.Fatalf("expected first data item to be object")
	}

	expectedKeys := []string{"id", "productID", "productName", "amount", "customerName", "status", "transactionDate", "createBy", "createOn"}
	for _, k := range expectedKeys {
		if _, ok := firstRow[k]; !ok {
			t.Errorf("E2: missing expected key %q in row", k)
		}
	}
}
