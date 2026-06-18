package repository

import (
	"database/sql"
)

type StatusRepo struct {
	db *sql.DB
}

func NewStatusRepo(db *sql.DB) *StatusRepo {
	return &StatusRepo{db: db}
}

func (r *StatusRepo) FindAll() ([]Status, error) {
	rows, err := r.db.Query(`SELECT id, name FROM statuses ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statuses []Status
	for rows.Next() {
		var s Status
		if err := rows.Scan(&s.ID, &s.Name); err != nil {
			return nil, err
		}
		statuses = append(statuses, s)
	}
	return statuses, rows.Err()
}
