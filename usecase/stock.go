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

func (uc GetLastPrice) Execute(ctx context.Context, symbol stock.Symbol) (float64, error) {
	info, err := uc.Integration.LastInfo(ctx, symbol)
	if err != nil {
		return 0, err
	}

	return info.LastPrice, err
}
