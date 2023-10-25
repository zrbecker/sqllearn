package chainstore_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/zrbecker/sqllearn/testutil"
	"github.com/zrbecker/sqllearn/v2_sqlc/gen/model"
	"github.com/zrbecker/sqllearn/v2_sqlc/schema"
	"github.com/zrbecker/sqllearn/v2_sqlc/stores/chainstore"
)

type ChainStoreSuite struct {
	testutil.PostgresSQLDockerSQLSuite
}

func (s *ChainStoreSuite) TestCreateChains() {
	db := s.DB()
	ctx := context.Background()

	store := chainstore.NewChainStore(db)

	s.Require().Nil(schema.CreateSchema(ctx, db))

	chains, err := store.CreateChains(ctx, []model.CreateChainParams{
		{Name: "Chain A", ChainID: "chain-a"},
		{Name: "Chain B", ChainID: "chain-b"},
	})
	s.Require().Nil(err)

	for i, chain := range chains {
		s.Require().NotZero(chain.ID)
		s.Require().NotZero(chain.CreatedAt)
		s.Require().NotZero(chain.UpdatedAt)

		chains[i].ID = *new(int32)
		chains[i].CreatedAt = *new(time.Time)
		chains[i].UpdatedAt = *new(time.Time)
	}

	s.Require().ElementsMatch([]model.Chain{
		{Name: "Chain A", ChainID: "chain-a"},
		{Name: "Chain B", ChainID: "chain-b"},
	}, chains)
}

func (s *ChainStoreSuite) TestGetChains() {
	db := s.DB()
	ctx := context.Background()

	store := chainstore.NewChainStore(db)

	s.Require().Nil(schema.CreateSchema(ctx, db))

	_, err := store.CreateChains(ctx, []model.CreateChainParams{
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

		chains[i].ID = *new(int32)
		chains[i].CreatedAt = *new(time.Time)
		chains[i].UpdatedAt = *new(time.Time)
	}

	s.Require().ElementsMatch([]model.Chain{
		{Name: "Chain A", ChainID: "chain-a"},
		{Name: "Chain B", ChainID: "chain-b"},
	}, chains)
}

func TestChainStoreSuite(t *testing.T) {
	suite.Run(t, &ChainStoreSuite{})
}
