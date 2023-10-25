package chainstore

import (
	"context"

	"github.com/zrbecker/sqllearn/v2_sqlc/gen/model"
)

type ChainStore struct {
	queries *model.Queries
}

func NewChainStore(db model.DBTX) *ChainStore {
	return &ChainStore{
		queries: model.New(db),
	}
}

func (s *ChainStore) CreateChains(
	ctx context.Context,
	params []model.CreateChainParams,
) ([]model.Chain, error) {
	var chains []model.Chain
	for _, param := range params {
		chain, err := s.queries.CreateChain(ctx, param)
		if err != nil {
			return nil, err
		}
		chains = append(chains, chain)
	}
	return chains, nil
}

func (s *ChainStore) GetChains(
	ctx context.Context,
	chainIDs []string,
) ([]model.Chain, error) {
	var chains []model.Chain
	for _, chainID := range chainIDs {
		chain, err := s.queries.GetChains(ctx, chainID)
		if err != nil {
			return nil, err
		}
		chains = append(chains, chain)
	}
	return chains, nil
}
