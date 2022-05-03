package stock

import (
	"context"
)

type (
	Symbol string

	Info struct {
		OpeningPrice float64
		MaxPrice     float64
		MinPrice     float64
		LastPrice    float64
		Change       float64
	}

	Provider interface {
		LastInfo(ctx context.Context, symbol Symbol) (Info, error)
	}
)
