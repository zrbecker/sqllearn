package chainstore_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/zrbecker/sqllearn/testutil"
	"github.com/zrbecker/sqllearn/v0_rawsql/model"
	"github.com/zrbecker/sqllearn/v0_rawsql/stores/chainstore"
)

type ChainStoreSuite struct {
	testutil.PostgresSQLDockerSQLSuite
}

func (s *ChainStoreSuite) TestCreateChains() {
	db := s.DB()
	store := chainstore.NewChainStore(db)
	ctx := context.Background()

	db.Exec(`
	CREATE TABLE chains (
		id BIGSERIAL PRIMARY KEY,
		created_at TIMESTAMPTZ,
		updated_at TIMESTAMPTZ,
		name TEXT NOT NULL,
		chain_id TEXT NOT NULL UNIQUE
	);
	`)

	chains, err := store.CreateChains(ctx, []chainstore.CreateChainDefinition{
		{Name: "Chain A", ChainID: "chain-a"},
		{Name: "Chain B", ChainID: "chain-b"},
	})
	s.Require().Nil(err)

	for i, chain := range chains {
		s.Require().NotZero(chain.ID)
		s.Require().NotZero(chain.CreatedAt)
		s.Require().NotZero(chain.UpdatedAt)

		chains[i].ID = *new(uint)
		chains[i].CreatedAt = *new(time.Time)
		chains[i].UpdatedAt = *new(time.Time)
	}

	s.Require().ElementsMatch([]model.Chain{
		{Name: "Chain A", ChainID: "chain-a"},
		{Name: "Chain B", ChainID: "chain-b"},
	}, chains)
}

func (s *ChainStoreSuite) TestQueryChains() {
	db := s.DB()
	store := chainstore.NewChainStore(db)
	ctx := context.Background()

	db.Exec(`
	CREATE TABLE chains (
		id BIGSERIAL PRIMARY KEY,
		created_at TIMESTAMPTZ,
		updated_at TIMESTAMPTZ,
		name TEXT NOT NULL,
		chain_id TEXT NOT NULL UNIQUE
	);
	`)

	_, err := store.CreateChains(ctx, []chainstore.CreateChainDefinition{
		{Name: "Chain A", ChainID: "chain-a"},
		{Name: "Chain B", ChainID: "chain-b"},
	})
	s.Require().Nil(err)

	chains, err := store.GetChains(ctx, []chainstore.GetChainDefinition{
		{Name: "Chain A", ChainID: "chain-a"},
		{Name: "Chain B", ChainID: "chain-b"},
	})
	s.Require().Nil(err)

	for i, chain := range chains {
		s.Require().NotZero(chain.ID)
		s.Require().NotZero(chain.CreatedAt)
		s.Require().NotZero(chain.UpdatedAt)

		chains[i].ID = *new(uint)
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
