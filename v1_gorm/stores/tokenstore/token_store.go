package tokenstore

import (
	"context"

	"github.com/zrbecker/sqllearn/v1_gorm/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TokenStore struct {
	db *gorm.DB
}

func NewTokenStore(db *gorm.DB) *TokenStore {
	return &TokenStore{db}
}

type CreateTokenDefintion struct {
	ChainID  string
	Name     string
	Denom    string
	Decimals int
}

func (s *TokenStore) CreateTokens(
	ctx context.Context,
	tokenDefs []CreateTokenDefintion,
) ([]model.Token, error) {
	db := s.db.WithContext(ctx)

	chainIDs := make([]string, 0)
	for _, tokenDef := range tokenDefs {
		chainIDs = append(chainIDs, tokenDef.ChainID)
	}

	var chains []model.Chain
	if err := db.Where("chain_id IN ?", chainIDs).Find(&chains).Error; err != nil {
		return nil, err
	}

	chainMap := make(map[string]model.Chain)
	for _, chain := range chains {
		chainMap[chain.ChainID] = chain
	}

	var tokens []model.Token
	for _, tokenDef := range tokenDefs {
		tokens = append(tokens, model.Token{
			ChainID:  chainMap[tokenDef.ChainID].ID,
			Name:     tokenDef.Name,
			Denom:    tokenDef.Denom,
			Decimals: tokenDef.Decimals,
		})
	}

	if err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "chain_id"}, {Name: "denom"}},
		UpdateAll: true,
	}).Create(&tokens).Error; err != nil {
		return nil, err
	}

	return tokens, nil
}

type GetTokenDefinition struct {
	ChainID string
	Denom   string
}

func (s *TokenStore) GetTokens(ctx context.Context, tokenDefs []GetTokenDefinition) ([]model.Token, error) {
	db := s.db.WithContext(ctx)

	chainIDs := make([]string, 0)
	for _, tokenDef := range tokenDefs {
		chainIDs = append(chainIDs, tokenDef.ChainID)
	}

	var chains []model.Chain
	if err := db.Where("chain_id IN ?", chainIDs).Find(&chains).Error; err != nil {
		return nil, err
	}

	chainMap := make(map[string]model.Chain)
	for _, chain := range chains {
		chainMap[chain.ChainID] = chain
	}

	var values [][]interface{}
	for _, tokenDef := range tokenDefs {
		values = append(values, []interface{}{
			chainMap[tokenDef.ChainID].ID,
			tokenDef.Denom,
		})
	}

	var tokens []model.Token
	if err := db.Where("(chain_id, denom) in ?", values).Find(&tokens).Error; err != nil {
		return nil, err
	}

	return tokens, nil
}
