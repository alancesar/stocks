package asset

import (
	"context"
	"fmt"
	"io"
	"stocks/currency"
	"stocks/separator"
	"stocks/stock"
)

type (
	Repository interface {
		Assets(ctx context.Context) (Assets, error)
	}

	Asset struct {
		Symbol       stock.Symbol
		Quantity     int
		AveragePrice currency.Currency
		LastPrice    currency.Currency
		Investment   currency.Currency
		Settled      currency.Currency
	}

	Assets []Asset
)

func (a Asset) Balance() currency.Currency {
	balance := a.Settled.Float64() - a.Investment.Float64()
	return currency.NewFromFloat(balance)
}

func (a Asset) GainLoss() currency.Currency {
	balance := float64(a.Quantity) * a.LastPrice.Float64()
	return currency.NewFromFloat(balance + a.Balance().Float64())
}

func (a Assets) Balance() currency.Currency {
	balance := 0.0

	for _, asset := range a {
		balance += asset.Balance().Float64()
	}

	return currency.NewFromFloat(balance)
}

func (a Assets) GainLoss() currency.Currency {
	gainLoss := 0.0

	for _, asset := range a {
		gainLoss += asset.GainLoss().Float64()
	}

	return currency.NewFromFloat(gainLoss)
}

func (a Assets) Print(writer io.Writer, sep separator.Separator) error {
	title := fmt.Sprintf("Symbol%sQtd.%sAvg. Price%sLast Price%sGain/Loss\n", sep, sep, sep, sep)
	if _, err := io.WriteString(writer, title); err != nil {
		return err
	}

	for _, asset := range a {
		line := fmt.Sprintf("%s%s%d%s%s%s%s%s%s\n",
			asset.Symbol, sep, asset.Quantity, sep, asset.AveragePrice, sep, asset.LastPrice, sep, asset.GainLoss())

		if _, err := io.WriteString(writer, line); err != nil {
			return err
		}
	}

	return nil
}
