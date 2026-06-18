package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PaulusBilly/internship-test-API/internal/dto"
	"github.com/PaulusBilly/internship-test-API/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockViewDataService struct {
	mock.Mock
}

func (m *MockViewDataService) BuildResponse() (*dto.ViewDataResponse, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ViewDataResponse), args.Error(1)
}

func TestViewDataHandler_GetViewData(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful response", func(t *testing.T) {
		mockService := new(MockViewDataService)
		expectedResponse := &dto.ViewDataResponse{
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
		}
		mockService.On("BuildResponse").Return(expectedResponse, nil)

		handler := NewViewDataHandler(mockService)
		router := gin.New()
		router.GET("/api/view-data", handler.GetViewData)

		req := httptest.NewRequest(http.MethodGet, "/api/view-data", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

		var actualResponse dto.ViewDataResponse
		err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse.Data[0].ID, actualResponse.Data[0].ID)
		assert.Equal(t, expectedResponse.Data[0].ProductID, actualResponse.Data[0].ProductID)
		assert.Equal(t, expectedResponse.Status[0].ID, actualResponse.Status[0].ID)
		assert.Equal(t, expectedResponse.Status[0].Name, actualResponse.Status[0].Name)

		mockService.AssertExpectations(t)
	})

	t.Run("service returns error", func(t *testing.T) {
		mockService := new(MockViewDataService)
		mockService.On("BuildResponse").Return(nil, assert.AnError)

		handler := NewViewDataHandler(mockService)
		router := gin.New()
		router.GET("/api/view-data", handler.GetViewData)

		req := httptest.NewRequest(http.MethodGet, "/api/view-data", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}
