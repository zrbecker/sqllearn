package chainstore_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zrbecker/sqllearn/testutil"
	"github.com/zrbecker/sqllearn/v1_gorm/model"
	"github.com/zrbecker/sqllearn/v1_gorm/stores/chainstore"
	"gorm.io/gorm"
)

type ChainStoreSuite struct {
	testutil.PostgresSQLDockerGORMSuite
}

func (s *ChainStoreSuite) TestCreateChains() {
	db := s.DB()
	ctx := context.Background()

	db.AutoMigrate(&model.Chain{})
	store := chainstore.NewChainStore(db)

	chains, err := store.CreateChains(ctx, []chainstore.CreateChainDefinition{
		{Name: "Chain A", ChainID: "chain-a"},
		{Name: "Chain B", ChainID: "chain-b"},
	})
	s.Require().Nil(err, "failed to create chains: %v", err)

	for i, chain := range chains {
		s.Require().NotZero(chain.ID)
		s.Require().NotZero(chain.CreatedAt)
		s.Require().NotZero(chain.UpdatedAt)
		chains[i].Model = *new(gorm.Model)
	}

	s.Require().ElementsMatch([]model.Chain{
		{Name: "Chain A", ChainID: "chain-a"},
		{Name: "Chain B", ChainID: "chain-b"},
	}, chains)
}

func (s *ChainStoreSuite) TestGetChains() {
	db := s.DB()
	ctx := context.Background()

	db.AutoMigrate(&model.Chain{})
	store := chainstore.NewChainStore(db)

	_, err := store.CreateChains(ctx, []chainstore.CreateChainDefinition{
		{Name: "Chain A", ChainID: "chain-a"},
		{Name: "Chain B", ChainID: "chain-b"},
	})
	s.Require().Nil(err)

	chains, err := store.GetChains(ctx, []string{"chain-a", "chain-b"})
	s.Require().Nil(err)

	for i, chain := range chains {
		s.Require().NotZero(chain.ID)
		s.Require().NotZero(chain.CreatedAt)
		s.Require().NotZero(chain.UpdatedAt)
		chains[i].Model = *new(gorm.Model)
	}

	s.Require().ElementsMatch([]model.Chain{
		{Name: "Chain A", ChainID: "chain-a"},
		{Name: "Chain B", ChainID: "chain-b"},
	}, chains)
}

func TestChainStoreSuite(t *testing.T) {
	suite.Run(t, &ChainStoreSuite{})
}
