package tokenstore

import (
	"context"

	"github.com/zrbecker/sqllearn/v2_sqlc/gen/model"
)

type TokenStore struct {
	queries *model.Queries
}

func NewTokenStore(db model.DBTX) *TokenStore {
	return &TokenStore{queries: model.New(db)}
}

type CreateTokenDefintion struct {
	ChainID  string
	Name     string
	Denom    string
	Decimals int32
}

func (s *TokenStore) CreateTokens(
	ctx context.Context,
	params []model.CreateTokenParams,
) ([]model.Token, error) {
	var tokens []model.Token
	for _, param := range params {
		token, err := s.queries.CreateToken(ctx, param)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}

func (s *TokenStore) GetTokens(
	ctx context.Context,
	params []model.GetTokenParams,
) ([]model.Token, error) {
	var tokens []model.Token
	for _, param := range params {
		token, err := s.queries.GetToken(ctx, param)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}
