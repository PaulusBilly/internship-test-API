package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"internship-test-API/internal/dto"
	"internship-test-API/internal/service"
)

type MockService struct {
	resp dto.ViewDataResponse
	err  error
}

func (m *MockService) BuildResponse() (dto.ViewDataResponse, error) {
	return m.resp, m.err
}

func TestGetViewData_HappyPath(t *testing.T) {
	mockService := &MockService{
		resp: dto.ViewDataResponse{
			Data: []dto.Row{
				{
					ID:              1372,
					ProductID:       "10001",
					ProductName:     "Test 1",
					Amount:          "1000",
					CustomerName:    "abc",
					Status:          0,
					TransactionDate: "2022-07-10 11:14:52",
					CreateBy:        "abc",
					CreateOn:        "2022-07-10 11:14:52",
				},
			},
			Status: []dto.StatusItem{
				{ID: 0, Name: "SUCCESS"},
				{ID: 1, Name: "FAILED"},
			},
		},
	}

	handler := &Handler{service: mockService}
	req := httptest.NewRequest("GET", "/api/view-data", nil)
	w := httptest.NewRecorder()

	handler.GetViewData(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	expectedContentType := "application/json"
	if ct := w.Header().Get("Content-Type"); ct != expectedContentType {
		t.Errorf("Expected Content-Type %s, got %s", expectedContentType, ct)
	}

	// Verify JSON contains expected structure
	body := w.Body.String()
	if !bytes.Contains([]byte(body), []byte(`"productID":"10001"`)) {
		t.Error("Response missing productID field")
	}
	if !bytes.Contains([]byte(body), []byte(`"amount":"1000"`)) {
		t.Error("Response missing amount field")
	}
	if !bytes.Contains([]byte(body), []byte(`"status":0`)) {
		t.Error("Response missing status field")
	}
}

func TestGetViewData_ServiceError(t *testing.T) {
	mockService := &MockService{
		err: service.ErrDatabase,
	}

	handler := &Handler{service: mockService}
	req := httptest.NewRequest("GET", "/api/view-data", nil)
	w := httptest.NewRecorder()

	handler.GetViewData(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}
