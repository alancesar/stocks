package stock

import (
	"context"
)

type (
	Symbol string

	Details struct {
		Symbol  Symbol
		Type    string
		Sector  string
		Segment string
	}

	Info struct {
		Symbol       Symbol
		OpeningPrice float64
		MaxPrice     float64
		MinPrice     float64
		LastPrice    float64
		Change       float64
	}

	Provider interface {
		Details(ctx context.Context, symbol Symbol) (Details, error)
		LastInfo(ctx context.Context, symbol Symbol) (Info, error)
	}
)
