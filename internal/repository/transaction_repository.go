package repository

import (
	"database/sql"

	"internship-test-api/internal/model"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) FindAll() ([]model.Transaction, error) {
	rows, err := r.db.Query(`
		SELECT t.id, t.product_id, t.product_name, t.amount_text, t.customer_name,
		       t.status_id, t.transaction_date, t.create_by, t.create_on
		FROM transactions t
		ORDER BY t.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []model.Transaction
	for rows.Next() {
		var t model.Transaction
		if err := rows.Scan(
			&t.ID,
			&t.ProductID,
			&t.ProductName,
			&t.AmountText,
			&t.CustomerName,
			&t.StatusID,
			&t.TransactionDate,
			&t.CreateBy,
			&t.CreateOn,
		); err != nil {
			return nil, err
		}
		txs = append(txs, t)
	}

	return txs, rows.Err()
}
