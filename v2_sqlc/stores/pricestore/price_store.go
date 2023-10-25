package pricestore

import (
	"context"

	"github.com/zrbecker/sqllearn/v2_sqlc/gen/model"
)

type PriceStore struct {
	queries *model.Queries
}

func NewPriceStore(db model.DBTX) *PriceStore {
	return &PriceStore{queries: model.New(db)}
}

func (s *PriceStore) CreatePrices(
	ctx context.Context,
	params []model.CreatePriceParams,
) ([]model.Price, error) {
	var prices []model.Price
	for _, param := range params {
		price, err := s.queries.CreatePrice(ctx, param)
		if err != nil {
			return nil, err
		}
		prices = append(prices, price)
	}
	return prices, nil
}

func (s *PriceStore) GetPrices(
	ctx context.Context,
	params []model.GetPricesParams,
) ([]model.Price, error) {
	var prices []model.Price
	for _, param := range params {
		tokenPrices, err := s.queries.GetPrices(ctx, param)
		if err != nil {
			return nil, err
		}
		prices = append(prices, tokenPrices...)
	}
	return prices, nil
}
