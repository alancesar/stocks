package operation

import (
	"context"
	"fmt"
	"io"
	"stocks/currency"
	"stocks/separator"
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
		List(ctx context.Context) (List, error)
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
		Stock     stock.Stock
		Type      Type
		Quantity  int
		UnitValue float64
		Date      time.Time
	}

	List []Operation
)

func (t Type) String() string {
	switch t {
	case Buy:
		return "BUY"
	case Sell:
		return "SELL"
	default:
		return ""
	}
}

func (o Operation) String() string {
	return fmt.Sprintf("Stock=%-6s Type=%s Quantity=%d UnitValue=%s Date=%s",
		o.Stock, o.Type, o.Quantity, currency.NewFromFloat(o.UnitValue), o.Date.Format("2006-01-02"))
}

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

func (l List) Print(writer io.Writer, sep separator.Separator) error {
	title := fmt.Sprintf("Stock%sType%sQtd%sUnit. Value%sDate\n", sep, sep, sep, sep)
	if _, err := io.WriteString(writer, title); err != nil {
		return err
	}

	for _, o := range l {
		line := fmt.Sprintf("%s%s%s%s%d%s%s%s%s\n",
			o.Stock, sep, o.Type, sep, o.Quantity, sep, currency.NewFromFloat(o.UnitValue), sep,
			o.Date.Format("2006-01-02"))

		if _, err := io.WriteString(writer, line); err != nil {
			return err
		}
	}

	return nil
}
