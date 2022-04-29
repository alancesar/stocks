package operation

import (
	"context"
	"stocks/stock"
	"time"
)

const (
	Buy  Type = "BUY"
	Sell Type = "SELL"
)

type (
	Type string

	Repository interface {
		Create(ctx context.Context, operation Operation) error
	}

	Operation struct {
		Type      Type
		Stock     stock.Stock
		Date      time.Time
		Quantity  int
		UnitValue float64
	}
)
