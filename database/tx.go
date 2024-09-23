package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewQueries(ctx context.Context, db *pgxpool.Pool) (*Queries, pgx.Tx, error) {
	tx, err := db.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}

	queries := New(tx)

	return queries, tx, nil
}
