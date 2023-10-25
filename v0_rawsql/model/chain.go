package model

import (
	"context"
	"time"

	"github.com/zrbecker/sqllearn/v0_rawsql/utils/sql"
)

type Chain struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	ChainID   string
}

func CreateSchema(ctx context.Context, db sql.DB) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE chains (
			id SERIAL PRIMARY KEY,
			created_at TIMESTAMPTZ,
			updated_at TIMESTAMPTZ,
		
			name TEXT NOT NULL,
			chain_id TEXT NOT NULL UNIQUE
		);
	`)
	return err
}
