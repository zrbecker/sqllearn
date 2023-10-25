package tokenstore_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/zrbecker/sqllearn/testutil"
	"github.com/zrbecker/sqllearn/v2_sqlc/gen/model"
	"github.com/zrbecker/sqllearn/v2_sqlc/schema"
	"github.com/zrbecker/sqllearn/v2_sqlc/stores/chainstore"
	"github.com/zrbecker/sqllearn/v2_sqlc/stores/tokenstore"
)

type TokenStoreSuite struct {
	testutil.PostgresSQLDockerSQLSuite
}

func (s *TokenStoreSuite) TestCreateTokens() {
	db := s.DB()
	ctx := context.Background()

	schema.CreateSchema(ctx, db)

	chainStore := chainstore.NewChainStore(db)
	store := tokenstore.NewTokenStore(db)

	chains, err := chainStore.CreateChains(ctx, []model.CreateChainParams{
		{Name: "Chain A", ChainID: "chain-a"},
		{Name: "Chain B", ChainID: "chain-b"},
	})
	s.Require().Nil(err)

	chainMap := make(map[string]model.Chain)
	for _, chain := range chains {
		chainMap[chain.ChainID] = chain
	}

	tokens, err := store.CreateTokens(ctx, []model.CreateTokenParams{
		{ChainID: "chain-a", Name: "Token A", Denom: "utokena", Decimals: 6},
		{ChainID: "chain-a", Name: "Token B", Denom: "utokenb", Decimals: 18},
		{ChainID: "chain-b", Name: "Token C", Denom: "utokenc", Decimals: 6},
	})
	s.Require().Nil(err)

	for i, token := range tokens {
		s.Require().NotZero(token.ID)
		s.Require().NotZero(token.CreatedAt)
		s.Require().NotZero(token.UpdatedAt)

		tokens[i].ID = *new(int32)
		tokens[i].CreatedAt = *new(time.Time)
		tokens[i].UpdatedAt = *new(time.Time)
	}

	s.Require().ElementsMatch([]model.Token{
		{ChainID: chainMap["chain-a"].ID, Name: "Token A", Denom: "utokena", Decimals: 6},
		{ChainID: chainMap["chain-a"].ID, Name: "Token B", Denom: "utokenb", Decimals: 18},
		{ChainID: chainMap["chain-b"].ID, Name: "Token C", Denom: "utokenc", Decimals: 6},
	}, tokens)
}

func (s *TokenStoreSuite) TestGetTokens() {
	db := s.DB()
	ctx := context.Background()

	schema.CreateSchema(ctx, db)

	chainStore := chainstore.NewChainStore(db)
	store := tokenstore.NewTokenStore(db)

	chains, err := chainStore.CreateChains(ctx, []model.CreateChainParams{
		{Name: "Chain A", ChainID: "chain-a"},
		{Name: "Chain B", ChainID: "chain-b"},
	})
	s.Require().Nil(err)

	chainMap := make(map[string]model.Chain)
	for _, chain := range chains {
		chainMap[chain.ChainID] = chain
	}

	_, err = store.CreateTokens(ctx, []model.CreateTokenParams{
		{ChainID: "chain-a", Name: "Token A", Denom: "utokena", Decimals: 6},
		{ChainID: "chain-a", Name: "Token B", Denom: "utokenb", Decimals: 18},
		{ChainID: "chain-b", Name: "Token C", Denom: "utokenc", Decimals: 6},
	})
	s.Require().Nil(err)

	tokens, err := store.GetTokens(ctx, []model.GetTokenParams{
		{ChainID: "chain-a", Denom: "utokena"},
		{ChainID: "chain-a", Denom: "utokenb"},
		{ChainID: "chain-b", Denom: "utokenc"},
	})
	s.Require().Nil(err)

	for i, token := range tokens {
		s.Require().NotZero(token.ID)
		s.Require().NotZero(token.CreatedAt)
		s.Require().NotZero(token.UpdatedAt)

		tokens[i].ID = *new(int32)
		tokens[i].CreatedAt = *new(time.Time)
		tokens[i].UpdatedAt = *new(time.Time)
	}

	s.Require().ElementsMatch([]model.Token{
		{ChainID: chainMap["chain-a"].ID, Name: "Token A", Denom: "utokena", Decimals: 6},
		{ChainID: chainMap["chain-a"].ID, Name: "Token B", Denom: "utokenb", Decimals: 18},
		{ChainID: chainMap["chain-b"].ID, Name: "Token C", Denom: "utokenc", Decimals: 6},
	}, tokens)
}

func TestTokenStoreSuite(t *testing.T) {
	suite.Run(t, &TokenStoreSuite{})
}
