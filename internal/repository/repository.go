package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/inasknh/simple-poke-app/internal/model"
)

type repository struct {
	db *sql.DB
}

// NewRepository creates a new instance of Repository.
func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

//go:generate mockery --name Repository --output ../mocks
type Repository interface {
	CreateBerry(ctx context.Context, berries []model.Berry) error
}

func (r *repository) CreateBerry(ctx context.Context, berries []model.Berry) error {
	if len(berries) == 0 {
		return nil
	}

	// Build query with placeholders (?, ?, ?)
	query := "INSERT INTO berries (name, url) VALUES "
	vals := []interface{}{}

	for i, u := range berries {
		if i > 0 {
			query += ","
		}
		query += "(?, ?)"
		vals = append(vals, u.Name, u.URL)
	}

	// Execute query
	_, err := r.db.Exec(query, vals...)
	if err != nil {
		return fmt.Errorf("failed to insert berries: %w", err)
	}

	return nil
}
