package repository

import (
	"time"
)

type Status struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Transaction struct {
	ID              int
	ProductID       string
	ProductName     string
	AmountText      string
	CustomerName    string
	StatusID        int
	TransactionDate time.Time
	CreateBy        string
	CreateOn        time.Time
}
