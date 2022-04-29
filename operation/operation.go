package operation

import (
	"context"
	"stocks/currency"
	"stocks/stock"
	"time"
)

const (
	Buy Type = iota
	Sell
)

type (
	Type int

	Repository interface {
		Create(ctx context.Context, operation Operation) error
	}

	ReportRepository interface {
		Summary(ctx context.Context) (Summary, error)
	}

	Entry struct {
		Stock        stock.Stock
		Quantity     int
		AveragePrice currency.Currency
		LastPrice    currency.Currency
		Investment   currency.Currency
		Settled      currency.Currency
	}

	Summary []Entry

	Report struct {
		Summary Summary
	}

	Operation struct {
		Type      Type
		Stock     stock.Stock
		Date      time.Time
		Quantity  int
		UnitValue float64
	}
)

func (s Entry) Balance() currency.Currency {
	balance := s.Settled.Float64() - s.Investment.Float64()
	return currency.NewFromFloat(balance)
}

func (s Entry) GainLoss() currency.Currency {
	balance := float64(s.Quantity) * s.LastPrice.Float64()
	return currency.NewFromFloat(balance + s.Balance().Float64())
}

func (r Report) Balance() currency.Currency {
	balance := 0.0

	for _, summary := range r.Summary {
		balance += summary.Balance().Float64()
	}

	return currency.NewFromFloat(balance)
}

func (r Report) GainLoss() currency.Currency {
	gainLoss := 0.0

	for _, summary := range r.Summary {
		gainLoss += summary.GainLoss().Float64()
	}

	return currency.NewFromFloat(gainLoss)
}

func (r Report) Print(writer io.Writer, sep separator.Separator) error {
	title := fmt.Sprintf("Stock%sQtd.%sAvg. Price%sLast Price%sGain/Loss\n", sep, sep, sep, sep)
	if _, err := io.WriteString(writer, title); err != nil {
		return err
	}

	for _, e := range r.Summary {
		line := fmt.Sprintf("%s%s%d%s%s%s%s%s%s\n",
			e.Stock, sep, e.Quantity, sep, e.AveragePrice, sep, e.LastPrice, sep, e.GainLoss())

		if _, err := io.WriteString(writer, line); err != nil {
			return err
		}
	}

	return nil
}
