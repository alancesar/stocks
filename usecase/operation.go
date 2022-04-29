package usecase

import (
	"context"
	"stocks/operation"
	"stocks/stock"
	"time"
)

type (
	BuyRequest struct {
		Stock     stock.Stock
		Quantity  int
		UnitValue float64
		Date      time.Time
	}

	BuyOperationUseCase struct {
		Repository operation.Repository
	}
)

func NewBuyOperationUseCase(repository operation.Repository) *BuyOperationUseCase {
	return &BuyOperationUseCase{
		Repository: repository,
	}
}

func (c BuyOperationUseCase) Execute(ctx context.Context, request BuyRequest) (operation.Operation, error) {
	op := operation.Operation{
		Type:      operation.Buy,
		Stock:     request.Stock,
		Date:      request.Date,
		Quantity:  request.Quantity,
		UnitValue: request.UnitValue,
	}

	if err := c.Repository.Create(ctx, op); err != nil {
		return operation.Operation{}, err
	}

	return op, nil
}
