package model

import "time"

type Transaction struct {
	ID               int64     `db:"id"`
	ProductID        string    `db:"product_id"`
	ProductName      string    `db:"product_name"`
	AmountText       string    `db:"amount_text"`
	CustomerName     string    `db:"customer_name"`
	StatusID         int       `db:"status_id"`
	TransactionDate  time.Time `db:"transaction_date"`
	CreateBy         string    `db:"create_by"`
	CreateOn         time.Time `db:"create_on"`
}
