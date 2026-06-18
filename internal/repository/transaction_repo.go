package repository

import (
	"database/sql"
	"time"
)

type TransactionRepo struct {
	db *sql.DB
}

func NewTransactionRepo(db *sql.DB) *TransactionRepo {
	return &TransactionRepo{db: db}
}

func (r *TransactionRepo) FindAll() ([]Transaction, error) {
	rows, err := r.db.Query(`
		SELECT 
			t.id, t.product_id, t.product_name, t.amount_text, t.customer_name,
			t.status_id, t.transaction_date, t.create_by, t.create_on
		FROM transactions t
		ORDER BY t.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []Transaction
	for rows.Next() {
		var t Transaction
		var txnDate, createOn sql.NullTime
		if err := rows.Scan(
			&t.ID, &t.ProductID, &t.ProductName, &t.AmountText,
			&t.CustomerName, &t.StatusID, &txnDate, &t.CreateBy, &createOn,
		); err != nil {
			return nil, err
		}
		if txnDate.Valid {
			t.TransactionDate = txnDate.Time
		}
		if createOn.Valid {
			t.CreateOn = createOn.Time
		}
		txs = append(txs, t)
	}
	return txs, rows.Err()
}
