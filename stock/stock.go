package stock

import (
	"context"
)

type (
	Stock string

	Info struct {
		OpeningPrice float64
		MaxPrice     float64
		MinPrice     float64
		LastPrice    float64
		Change       float64
	}

	Provider interface {
		LastInfo(ctx context.Context, stock Stock) (Info, error)
	}
)
