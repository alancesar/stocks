package stock

import (
	"context"
	"time"
)

type (
	Stock string

	Info struct {
		ClosingDate  time.Time
		OpeningPrice float64
		MaxPrice     float64
		MinPrice     float64
		MediumPrice  float64
		LastPrice    float64
	}

	Provider interface {
		LastInfo(ctx context.Context, stock Stock) (Info, error)
	}
)
