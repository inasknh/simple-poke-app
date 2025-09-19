package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/inasknh/simple-poke-app/internal/model"
)

const (
	getAllBerries = "SELECT name, url FROM berries"
)

type repository struct {
	db *sql.DB
}

// NewRepository creates a new instance of Repository.
func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

type Repository interface {
	CreateBerry(ctx context.Context, berries []model.Berry) error
	FetchBerries(ctx context.Context) (*model.BerriesResponse, error)
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
	_, err := r.db.ExecContext(ctx, query, vals...)
	if err != nil {
		return fmt.Errorf("failed to insert berries: %w", err)
	}

	return nil
}

func (r *repository) FetchBerries(ctx context.Context) (*model.BerriesResponse, error) {
	rows, err := r.db.QueryContext(ctx, getAllBerries)
	if err != nil {
		return nil, err
	}

	res := []model.Berry{}
	for rows.Next() {
		var b model.Berry
		err = rows.Scan(
			&b.Name,
			&b.URL,
		)

		if err != nil {
			return nil, err
		}

		res = append(res, b)
	}

	return &model.BerriesResponse{Berries: res}, nil
}
