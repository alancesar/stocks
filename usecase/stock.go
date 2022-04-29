package usecase

import (
	"context"
	"stocks/stock"
)

type (
	GetLastPrice struct {
		Integration stock.Provider
	}
)

func NewGetLastPrice(stockService stock.Provider) *GetLastPrice {
	return &GetLastPrice{
		Integration: stockService,
	}
}

func (p GetLastPrice) Execute(ctx context.Context, s stock.Stock) (float64, error) {
	info, err := p.Integration.LastInfo(ctx, s)
	if err != nil {
		return 0, err
	}

	return info.LastPrice, err
}
