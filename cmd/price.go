package main

import (
	"errors"
	"stocks/stock"
)

func CreatePriceRequest(args ...string) (stock.Stock, error) {
	if len(args) < 1 || args[0] == "" {
		return "", errors.New("usage: stocks price <symbol>")
	}

	return stock.Stock(args[0]), nil
}
