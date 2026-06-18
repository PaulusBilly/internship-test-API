package models

import "time"

type Status struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Transaction struct {
	ID               int64     `json:"-"`
	ProductID        string    `json:"-"`
	ProductName      string    `json:"-"`
	AmountText       string    `json:"-"`
	CustomerName     string    `json:"-"`
	StatusID         int       `json:"-"`
	TransactionDate  time.Time `json:"-"`
	CreateBy         string    `json:"-"`
	CreateOn         time.Time `json:"-"`
}
