package repository

import (
	"database/sql"

	"internship-test-api/internal/model"
)

type StatusRepository struct {
	db *sql.DB
}

func NewStatusRepository(db *sql.DB) *StatusRepository {
	return &StatusRepository{db: db}
}

func (r *StatusRepository) FindAll() ([]model.Status, error) {
	rows, err := r.db.Query(`SELECT id, name FROM statuses ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statuses []model.Status
	for rows.Next() {
		var s model.Status
		if err := rows.Scan(&s.ID, &s.Name); err != nil {
			return nil, err
		}
		statuses = append(statuses, s)
	}

	return statuses, rows.Err()
}
