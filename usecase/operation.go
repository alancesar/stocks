package usecase

import (
	"context"
	"stocks/asset"
	"stocks/currency"
	"stocks/operation"
	"stocks/stock"
	"sync"
	"time"
)

type (
	Fetcher interface {
		Fetch(ctx context.Context, symbol stock.Symbol) error
	}

	BuyRequest struct {
		Symbol    stock.Symbol
		Quantity  int
		UnitValue float64
		Date      time.Time
	}

	BuyOperationUseCase struct {
		Fetcher    Fetcher
		Repository operation.Repository
	}

	ListUseCase struct {
		Repository operation.Repository
	}

	AssetsUseCase struct {
		Provider   stock.Provider
		Repository asset.Repository
	}
)

func NewBuyOperationUseCase(repository operation.Repository, fetcher Fetcher) *BuyOperationUseCase {
	return &BuyOperationUseCase{
		Repository: repository,
		Fetcher:    fetcher,
	}
}

func NewListUseCase(repository operation.Repository) *ListUseCase {
	return &ListUseCase{
		Repository: repository,
	}
}

func NewAssetsUseCase(provider stock.Provider, repository asset.Repository) *AssetsUseCase {
	return &AssetsUseCase{
		Provider:   provider,
		Repository: repository,
	}
}

func (uc BuyOperationUseCase) Execute(ctx context.Context, request BuyRequest) (operation.Operation, error) {
	if err := uc.Fetcher.Fetch(ctx, request.Symbol); err != nil {
		return operation.Operation{}, err
	}

	op := operation.Operation{
		Type:      operation.Buy,
		Symbol:    request.Symbol,
		Date:      request.Date,
		Quantity:  request.Quantity,
		UnitValue: request.UnitValue,
	}

	if err := uc.Repository.Create(ctx, op); err != nil {
		return operation.Operation{}, err
	}

	return op, nil
}

func (uc ListUseCase) Execute(ctx context.Context) (operation.List, error) {
	return uc.Repository.List(ctx)
}

func (uc AssetsUseCase) Execute(ctx context.Context) (asset.Assets, error) {
	assets, err := uc.Repository.Assets(ctx)
	if err != nil {
		return nil, err
	}

	done := make(chan bool)
	fail := make(chan error)

	wg := sync.WaitGroup{}
	for i := range assets {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			info, err := uc.Provider.LastInfo(ctx, assets[i].Symbol)
			if err != nil {
				fail <- err
				return
			}

			assets[i].LastPrice = currency.NewFromFloat(info.LastPrice)
		}(i)
	}

	go func() {
		wg.Wait()
		done <- true
	}()

	select {
	case <-done:
		return assets, nil
	case err := <-fail:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
