package schema

import (
	"context"

	"github.com/zrbecker/sqllearn/v2_sqlc/gen/model"
)

func CreateSchema(ctx context.Context, db model.DBTX) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE chains (
			id SERIAL PRIMARY KEY,
			created_at TIMESTAMPTZ NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL,
		
			name TEXT NOT NULL,
			chain_id TEXT NOT NULL UNIQUE
		);
		
		CREATE TABLE tokens (
			id SERIAL PRIMARY KEY,
			created_at TIMESTAMPTZ NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL,
		
			chain_id INT NOT NULL,
			name TEXT NOT NULL,
			denom TEXT NOT NULL,
			decimals INT NOT NULL,
		
			CONSTRAINT fk_tokens_chain
				FOREIGN KEY (chain_id)
					REFERENCES chains(id)
					ON DELETE CASCADE
		);
		
		CREATE UNIQUE INDEX ux_tokens_chain_id_denom ON tokens(chain_id, denom);
		
		CREATE TABLE prices (
			id SERIAL PRIMARY KEY,
			created_at TIMESTAMPTZ NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL,
		
			token_id INT NOT NULL,
			timestamp TIMESTAMPTZ NOT NULL,
			price TEXT NOT NULL
		);
		
		CREATE UNIQUE INDEX ux_prices_token_id_timestamp ON prices(token_id, timestamp);
	`)
	return err
}
