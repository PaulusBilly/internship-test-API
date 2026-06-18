package service

import (
	"fmt"
	"time"

	"github.com/PaulusBilly/internship-test-API/internal/dto"
	"github.com/PaulusBilly/internship-test-API/internal/models"
	"github.com/PaulusBilly/internship-test-API/internal/repository"
)

const dateLayout = "2006-01-02 15:04:05"

type Service struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) BuildResponse() (dto.ViewDataResponse, error) {
	txs, err := s.repo.GetTransactions()
	if err != nil {
		return dto.ViewDataResponse{}, err
	}
	statuses, err := s.repo.GetStatuses()
	if err != nil {
		return dto.ViewDataResponse{}, err
	}

	var data []dto.Row
	for _, t := range txs {
		row := dto.Row{
			ID:              t.ID,
			ProductID:       t.ProductID,
			ProductName:     t.ProductName,
			Amount:          t.AmountText,
			CustomerName:    t.CustomerName,
			Status:          t.StatusID,
			TransactionDate: t.TransactionDate.Format(dateLayout),
			CreateBy:        t.CreateBy,
			CreateOn:        t.CreateOn.Format(dateLayout),
		}
		data = append(data, row)
	}

	var statusList []dto.StatusItem
	for _, st := range statuses {
		statusList = append(statusList, dto.StatusItem{
			ID:   st.ID,
			Name: st.Name,
		})
	}

	return dto.ViewDataResponse{
		Data:   data,
		Status: statusList,
	}, nil
}

func (s *Service) Seed() error {
	count, err := s.repo.CountTransactions()
	if err != nil {
		return fmt.Errorf("failed to count transactions: %w", err)
	}
	if count > 0 {
		return nil // already seeded
	}

	return nil // placeholder: actual seeding handled in bootstrap
}
