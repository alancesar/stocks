package stock

import (
	"context"
)

type (
	Symbol string

	Fetcher struct {
		Repository Repository
		Provider   Provider
	}

	Details struct {
		Symbol    Symbol
		Name      string
		Sector    string
		SubSector string
		Segment   string
	}

	Info struct {
		Symbol       Symbol
		OpeningPrice float64
		MaxPrice     float64
		MinPrice     float64
		LastPrice    float64
		Change       float64
	}

	Repository interface {
		GetDetails(ctx context.Context, symbol Symbol) (Details, error)
		InsertDetails(ctx context.Context, details Details) error
	}

	Provider interface {
		Details(ctx context.Context, symbol Symbol) (Details, error)
		LastInfo(ctx context.Context, symbol Symbol) (Info, error)
	}
)

func NewFetcher(repository Repository, provider Provider) *Fetcher {
	return &Fetcher{
		Repository: repository,
		Provider:   provider,
	}
}

func (f Fetcher) Fetch(ctx context.Context, symbol Symbol) error {
	if _, err := f.Repository.GetDetails(ctx, symbol); err == nil {
		return nil
	}

	if details, err := f.Provider.Details(ctx, symbol); err != nil {
		return err
	} else {
		return f.Repository.InsertDetails(ctx, details)
	}
}
