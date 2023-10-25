package chainstore

import (
	"context"

	"github.com/zrbecker/sqllearn/v1_gorm/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ChainStore struct {
	db *gorm.DB
}

func NewChainStore(db *gorm.DB) *ChainStore {
	return &ChainStore{db}
}

type CreateChainDefinition struct {
	Name    string
	ChainID string
}

func (s *ChainStore) CreateChains(
	ctx context.Context,
	chainDefs []CreateChainDefinition,
) ([]model.Chain, error) {
	db := s.db.WithContext(ctx)

	var chains []model.Chain
	for _, chainDef := range chainDefs {
		chains = append(chains, model.Chain{Name: chainDef.Name, ChainID: chainDef.ChainID})
	}

	if err := db.
		Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "chain_id"}}, UpdateAll: true}).
		Create(&chains).Error; err != nil {
		return nil, err
	}

	return chains, nil
}

func (s *ChainStore) GetChains(ctx context.Context, chainIDs []string) ([]model.Chain, error) {
	db := s.db.WithContext(ctx)

	var chains []model.Chain
	if tx := db.Where("chain_id IN ?", chainIDs).Find(&chains); tx.Error != nil {
		return nil, tx.Error
	}

	return chains, nil
}
