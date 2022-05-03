package usecase

import (
	"context"
	"stocks/currency"
	"stocks/operation"
	"stocks/stock"
	"sync"
	"time"
)

type (
	BuyRequest struct {
		Symbol    stock.Symbol
		Quantity  int
		UnitValue float64
		Date      time.Time
	}

	BuyOperationUseCase struct {
		Repository operation.Repository
	}

	ListUseCase struct {
		Repository operation.Repository
	}

	ReportUseCase struct {
		Provider   stock.Provider
		Repository operation.ReportRepository
	}
)

func NewBuyOperationUseCase(repository operation.Repository) *BuyOperationUseCase {
	return &BuyOperationUseCase{
		Repository: repository,
	}
}

func NewListUseCase(repository operation.Repository) *ListUseCase {
	return &ListUseCase{
		Repository: repository,
	}
}

func NewReportUseCase(provider stock.Provider, repository operation.ReportRepository) *ReportUseCase {
	return &ReportUseCase{
		Provider:   provider,
		Repository: repository,
	}
}

func (uc BuyOperationUseCase) Execute(ctx context.Context, request BuyRequest) (operation.Operation, error) {
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

func (uc ReportUseCase) Execute(ctx context.Context) (operation.Report, error) {
	summary, err := uc.Repository.Summary(ctx)
	if err != nil {
		return operation.Report{}, err
	}

	done := make(chan bool)
	fail := make(chan error)

	wg := sync.WaitGroup{}
	for i := range summary {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			info, err := uc.Provider.LastInfo(ctx, summary[i].Symbol)
			if err != nil {
				fail <- err
				return
			}

			summary[i].LastPrice = currency.NewFromFloat(info.LastPrice)
		}(i)
	}

	go func() {
		wg.Wait()
		done <- true
	}()

	select {
	case <-done:
		return operation.Report{
			Summary: summary,
		}, nil
	case err := <-fail:
		return operation.Report{}, err
	case <-ctx.Done():
		return operation.Report{}, ctx.Err()
	}
}
