package pricestore_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/zrbecker/sqllearn/testutil"
	"github.com/zrbecker/sqllearn/v2_sqlc/gen/model"
	"github.com/zrbecker/sqllearn/v2_sqlc/schema"
	"github.com/zrbecker/sqllearn/v2_sqlc/stores/chainstore"
	"github.com/zrbecker/sqllearn/v2_sqlc/stores/pricestore"
	"github.com/zrbecker/sqllearn/v2_sqlc/stores/tokenstore"
)

func parseTimeNow(value string) time.Time {
	t, err := time.Parse(time.DateTime, value)
	if err != nil {
		panic(err)
	}
	return t
}

type PriceStoreSuite struct {
	testutil.PostgresSQLDockerSQLSuite
}

func (s *PriceStoreSuite) TestCreatePrices() {
	db := s.DB()
	ctx := context.Background()

	schema.CreateSchema(ctx, db)

	chainStore := chainstore.NewChainStore(db)
	tokenStore := tokenstore.NewTokenStore(db)
	store := pricestore.NewPriceStore(db)

	_, err := chainStore.CreateChains(ctx, []model.CreateChainParams{
		{Name: "Chain A", ChainID: "chain-a"},
	})
	s.Require().Nil(err)

	tokens, err := tokenStore.CreateTokens(ctx, []model.CreateTokenParams{
		{ChainID: "chain-a", Name: "Token A", Denom: "utokena", Decimals: 6},
	})
	s.Require().Nil(err)
	token := tokens[0]

	prices, err := store.CreatePrices(ctx, []model.CreatePriceParams{{
		ChainID:   "chain-a",
		Denom:     "utokena",
		Price:     "4.20",
		Timestamp: parseTimeNow("2023-06-01 12:00:00"),
	}, {
		ChainID:   "chain-a",
		Denom:     "utokena",
		Price:     "4.30",
		Timestamp: parseTimeNow("2023-06-02 12:00:00"),
	}, {
		ChainID:   "chain-a",
		Denom:     "utokena",
		Price:     "4.40",
		Timestamp: parseTimeNow("2023-06-03 12:00:00"),
	}})
	s.Require().Nil(err)

	for i, price := range prices {
		s.Require().NotZero(price.ID)
		s.Require().NotZero(price.CreatedAt)
		s.Require().NotZero(price.UpdatedAt)

		prices[i].ID = *new(int32)
		prices[i].CreatedAt = *new(time.Time)
		prices[i].UpdatedAt = *new(time.Time)
	}

	for i := range prices {
		prices[i].Timestamp = prices[i].Timestamp.UTC()
	}

	s.Require().ElementsMatch([]model.Price{
		{TokenID: token.ID, Timestamp: parseTimeNow("2023-06-01 12:00:00"), Price: "4.20"},
		{TokenID: token.ID, Timestamp: parseTimeNow("2023-06-02 12:00:00"), Price: "4.30"},
		{TokenID: token.ID, Timestamp: parseTimeNow("2023-06-03 12:00:00"), Price: "4.40"},
	}, prices)
}

func (s *PriceStoreSuite) TestGetPrices() {
	db := s.DB()
	ctx := context.Background()

	schema.CreateSchema(ctx, db)

	chainStore := chainstore.NewChainStore(db)
	tokenStore := tokenstore.NewTokenStore(db)
	store := pricestore.NewPriceStore(db)

	_, err := chainStore.CreateChains(ctx, []model.CreateChainParams{
		{Name: "Chain A", ChainID: "chain-a"},
	})
	s.Require().Nil(err)

	tokens, err := tokenStore.CreateTokens(ctx, []model.CreateTokenParams{
		{ChainID: "chain-a", Name: "Token A", Denom: "utokena", Decimals: 6},
	})
	s.Require().Nil(err)
	token := tokens[0]

	_, err = store.CreatePrices(ctx, []model.CreatePriceParams{{
		ChainID:   "chain-a",
		Denom:     "utokena",
		Price:     "4.20",
		Timestamp: parseTimeNow("2023-06-01 12:00:00"),
	}, {
		ChainID:   "chain-a",
		Denom:     "utokena",
		Price:     "4.30",
		Timestamp: parseTimeNow("2023-06-02 12:00:00"),
	}, {
		ChainID:   "chain-a",
		Denom:     "utokena",
		Price:     "4.40",
		Timestamp: parseTimeNow("2023-06-03 12:00:00"),
	}})
	s.Require().Nil(err)

	prices, err := store.GetPrices(ctx, []model.GetPricesParams{{
		ChainID: "chain-a",
		Denom:   "utokena",
	}})
	s.Require().Nil(err)

	for i, price := range prices {
		s.Require().NotZero(price.ID)
		s.Require().NotZero(price.CreatedAt)
		s.Require().NotZero(price.UpdatedAt)

		prices[i].ID = *new(int32)
		prices[i].CreatedAt = *new(time.Time)
		prices[i].UpdatedAt = *new(time.Time)
	}

	for i := range prices {
		prices[i].Timestamp = prices[i].Timestamp.UTC()
	}

	s.Require().ElementsMatch([]model.Price{
		{TokenID: token.ID, Timestamp: parseTimeNow("2023-06-01 12:00:00"), Price: "4.20"},
		{TokenID: token.ID, Timestamp: parseTimeNow("2023-06-02 12:00:00"), Price: "4.30"},
		{TokenID: token.ID, Timestamp: parseTimeNow("2023-06-03 12:00:00"), Price: "4.40"},
	}, prices)
}

func TestPriceStoreSuite(t *testing.T) {
	suite.Run(t, &PriceStoreSuite{})
}
