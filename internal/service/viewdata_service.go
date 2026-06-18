package service

import (
	"fmt"
	"time"

	"internship-test-api/internal/handler"
	"internship-test-api/internal/model"
	"internship-test-api/internal/repository"
)

type ViewDataService struct {
	statusRepo      *repository.StatusRepository
	transactionRepo *repository.TransactionRepository
}

func NewViewDataService(statusRepo *repository.StatusRepository, transactionRepo *repository.TransactionRepository) *ViewDataService {
	return &ViewDataService{
		statusRepo:      statusRepo,
		transactionRepo: transactionRepo,
	}
}

func (s *ViewDataService) GetViewData() (*handler.ViewDataResponse, error) {
	txs, err := s.transactionRepo.FindAll()
	if err != nil {
		return nil, err
	}

	statuses, err := s.statusRepo.FindAll()
	if err != nil {
		return nil, err
	}

	response := &handler.ViewDataResponse{
		Data:   make([]handler.Row, 0, len(txs)),
		Status: make([]handler.StatusItem, 0, len(statuses)),
	}

	for _, tx := range txs {
		row := handler.Row{
			ID:              tx.ID,
			ProductID:       tx.ProductID,
			ProductName:     tx.ProductName,
			Amount:          tx.AmountText,
			CustomerName:    tx.CustomerName,
			Status:          tx.StatusID,
			TransactionDate: tx.TransactionDate.Format("2006-01-02 15:04:05"),
			CreateBy:        tx.CreateBy,
			CreateOn:        tx.CreateOn.Format("2006-01-02 15:04:05"),
		}
		response.Data = append(response.Data, row)
	}

	for _, st := range statuses {
		response.Status = append(response.Status, handler.StatusItem{
			ID:   st.ID,
			Name: st.Name,
		})
	}

	return response, nil
}
