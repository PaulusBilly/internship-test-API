package repository

import (
	"database/sql"
	"github.com/PaulusBilly/internship-test-API/internal/models"
)

type Repository interface {
	GetTransactions() ([]models.Transaction, error)
	GetStatuses() ([]models.Status, error)
	SaveStatus(s models.Status) error
	SaveTransaction(t models.Transaction) error
	CountTransactions() (int64, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) GetTransactions() ([]models.Transaction, error) {
	rows, err := r.db.Query(`
		SELECT id, product_id, product_name, amount_text, customer_name,
		       status_id, transaction_date, create_by, create_on
		FROM transactions`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []models.Transaction
	for rows.Next() {
		var t models.Transaction
		err := rows.Scan(
			&t.ID,
			&t.ProductID,
			&t.ProductName,
			&t.AmountText,
			&t.CustomerName,
			&t.StatusID,
			&t.TransactionDate,
			&t.CreateBy,
			&t.CreateOn,
		)
		if err != nil {
			return nil, err
		}
		txs = append(txs, t)
	}
	return txs, rows.Err()
}

func (r *postgresRepository) GetStatuses() ([]models.Status, error) {
	rows, err := r.db.Query(`SELECT id, name FROM statuses`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statuses []models.Status
	for rows.Next() {
		var s models.Status
		if err := rows.Scan(&s.ID, &s.Name); err != nil {
			return nil, err
		}
		statuses = append(statuses, s)
	}
	return statuses, rows.Err()
}

func (r *postgresRepository) SaveStatus(s models.Status) error {
	_, err := r.db.Exec(`INSERT INTO statuses (id, name) VALUES ($1, $2) ON CONFLICT (id) DO NOTHING`, s.ID, s.Name)
	return err
}

func (r *postgresRepository) SaveTransaction(t models.Transaction) error {
	_, err := r.db.Exec(`
		INSERT INTO transactions (
			id, product_id, product_name, amount_text, customer_name,
			status_id, transaction_date, create_by, create_on
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) DO NOTHING`,
		t.ID, t.ProductID, t.ProductName, t.AmountText, t.CustomerName,
		t.StatusID, t.TransactionDate, t.CreateBy, t.CreateOn,
	)
	return err
}

func (r *postgresRepository) CountTransactions() (int64, error) {
	var count int64
	err := r.db.QueryRow(`SELECT COUNT(*) FROM transactions`).Scan(&count)
	return count, err
}
