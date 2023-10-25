package tokenstore_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zrbecker/sqllearn/testutil"
	"github.com/zrbecker/sqllearn/v1_gorm/model"
	"github.com/zrbecker/sqllearn/v1_gorm/stores/chainstore"
	"github.com/zrbecker/sqllearn/v1_gorm/stores/tokenstore"
	"gorm.io/gorm"
)

type TokenStoreSuite struct {
	testutil.PostgresSQLDockerGORMSuite
}

func (s *TokenStoreSuite) TestCreateToken() {
	db := s.DB()
	ctx := context.Background()

	db.AutoMigrate(&model.Chain{}, &model.Token{})

	chainStore := chainstore.NewChainStore(db)
	store := tokenstore.NewTokenStore(db)

	chains, err := chainStore.CreateChains(ctx, []chainstore.CreateChainDefinition{
		{Name: "Chain A", ChainID: "chain-a"},
		{Name: "Chain B", ChainID: "chain-b"},
	})
	s.Require().Nil(err)

	chainMap := make(map[string]model.Chain)
	for _, chain := range chains {
		chainMap[chain.ChainID] = chain
	}

	tokens, err := store.CreateTokens(ctx, []tokenstore.CreateTokenDefintion{
		{ChainID: "chain-a", Name: "Token A", Denom: "utokena", Decimals: 6},
		{ChainID: "chain-a", Name: "Token B", Denom: "utokenb", Decimals: 18},
		{ChainID: "chain-b", Name: "Token C", Denom: "utokenc", Decimals: 6},
	})
	s.Require().Nil(err)

	for i := range tokens {
		tokens[i].Model = *new(gorm.Model)
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

	db.AutoMigrate(&model.Chain{}, &model.Token{})

	chainStore := chainstore.NewChainStore(db)
	store := tokenstore.NewTokenStore(db)

	chains, err := chainStore.CreateChains(ctx, []chainstore.CreateChainDefinition{
		{Name: "Chain A", ChainID: "chain-a"},
		{Name: "Chain B", ChainID: "chain-b"},
	})
	s.Require().Nil(err)

	chainMap := make(map[string]model.Chain)
	for _, chain := range chains {
		chainMap[chain.ChainID] = chain
	}

	_, err = store.CreateTokens(ctx, []tokenstore.CreateTokenDefintion{
		{ChainID: "chain-a", Name: "Token A", Denom: "utokena", Decimals: 6},
		{ChainID: "chain-a", Name: "Token B", Denom: "utokenb", Decimals: 18},
		{ChainID: "chain-b", Name: "Token C", Denom: "utokenc", Decimals: 6},
	})
	s.Require().Nil(err)

	tokens, err := store.GetTokens(ctx, []tokenstore.GetTokenDefinition{
		{ChainID: "chain-a", Denom: "utokena"},
		{ChainID: "chain-a", Denom: "utokenb"},
		{ChainID: "chain-b", Denom: "utokenc"},
	})
	s.Require().Nil(err)

	for i := range tokens {
		tokens[i].Model = *new(gorm.Model)
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
