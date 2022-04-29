package main

import (
	"errors"
	"stocks/date"
	"stocks/stock"
	"stocks/usecase"
	"strconv"
)

func CreateBuyRequest(args ...string) (usecase.BuyRequest, error) {
	if len(args) < 3 || len(args) > 4 {
		return usecase.BuyRequest{}, errors.New("usage: stocks buy <symbol> <quantity> <unit-value> [<date>]")
	}

	quantity, err := strconv.Atoi(args[1])
	if err != nil {
		return usecase.BuyRequest{}, errors.New("invalid quantity")
	}

	value, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return usecase.BuyRequest{}, errors.New("invalid value format")
	}

	rawDate := "today"
	if len(args) == 4 {
		rawDate = args[3]
	}

	d, err := date.Parse(rawDate)
	if err != nil {
		return usecase.BuyRequest{}, err
	}

	return usecase.BuyRequest{
		Stock:     stock.Stock(args[0]),
		Quantity:  quantity,
		UnitValue: value,
		Date:      d,
	}, nil
}
