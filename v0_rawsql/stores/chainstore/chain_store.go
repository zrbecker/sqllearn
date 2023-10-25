package chainstore

import (
	"context"
	"fmt"
	"strings"

	"github.com/zrbecker/sqllearn/v0_rawsql/model"
	"github.com/zrbecker/sqllearn/v0_rawsql/utils/sql"
)

type ChainStore struct {
	db sql.DB
}

func NewChainStore(db sql.DB) *ChainStore {
	return &ChainStore{db: db}
}

type CreateChainDefinition struct {
	Name    string
	ChainID string
}

func (s *ChainStore) CreateChains(
	ctx context.Context,
	chainDefs []CreateChainDefinition,
) ([]model.Chain, error) {
	builder := make([]string, 0)
	values := make([]any, 0)

	builder = append(builder, "INSERT INTO chains(created_at, updated_at, name, chain_id) VALUES")
	for i, chainDef := range chainDefs {
		builder = append(builder, fmt.Sprintf(" (NOW(), NOW(), $%d, $%d),", 2*i+1, 2*i+2))
		values = append(values, chainDef.Name)
		values = append(values, chainDef.ChainID)
	}
	builder[len(builder)-1] = strings.TrimSuffix(builder[len(builder)-1], ",")

	builder = append(builder, " ON CONFLICT (chain_id) DO UPDATE SET updated_at=NOW(), name=EXCLUDED.name")
	builder = append(builder, " RETURNING id, created_at, updated_at, name, chain_id;")
	query := strings.Join(builder, "")

	rows, err := s.db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, err
	}

	chains := make([]model.Chain, 0)
	for rows.Next() {
		chains = append(chains, model.Chain{})
		chain := &chains[len(chains)-1]
		if err := rows.Scan(
			&chain.ID,
			&chain.CreatedAt,
			&chain.UpdatedAt,
			&chain.Name,
			&chain.ChainID,
		); err != nil {
			return nil, err
		}
	}

	return chains, nil
}

type GetChainDefinition struct {
	Name    string
	ChainID string
}

func (s *ChainStore) GetChains(
	ctx context.Context,
	chainDefs []GetChainDefinition,
) ([]model.Chain, error) {
	builder := make([]string, 0)
	values := make([]any, 0)

	builder = append(builder, "SELECT id, created_at, updated_at, name, chain_id FROM chains WHERE")
	for i, chainDef := range chainDefs {
		builder = append(builder, fmt.Sprintf(" (name=$%d AND chain_id=$%d) OR", 2*i+1, 2*i+2))
		values = append(values, chainDef.Name)
		values = append(values, chainDef.ChainID)
	}
	builder[len(builder)-1] = strings.TrimSuffix(builder[len(builder)-1], " OR")
	builder = append(builder, " ORDER BY chain_id;")
	query := strings.Join(builder, "")

	rows, err := s.db.QueryContext(ctx, query, values...)
	if err != nil {
		return nil, err
	}

	chains := make([]model.Chain, 0)
	for rows.Next() {
		chains = append(chains, model.Chain{})
		chain := &chains[len(chains)-1]
		if err := rows.Scan(
			&chain.ID,
			&chain.CreatedAt,
			&chain.UpdatedAt,
			&chain.Name,
			&chain.ChainID,
		); err != nil {
			return nil, err
		}
	}

	return chains, nil
}
