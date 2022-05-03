package mfinance

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"stocks/stock"
)

const (
	baseUrl = "https://mfinance.com.br/api/v1/stocks"
)

type (
	Response struct {
		Change       float64 `json:"change"`
		ClosingPrice float64 `json:"closingPrice"`
		Eps          float64 `json:"eps"`
		High         float64 `json:"high"`
		LastPrice    float64 `json:"lastPrice"`
		LastYearHigh float64 `json:"lastYearHigh"`
		LastYearLow  float64 `json:"lastYearLow"`
		Low          float64 `json:"low"`
		MarketCap    int64   `json:"marketCap"`
		Name         string  `json:"name"`
		Pe           float64 `json:"pe"`
		PriceOpen    float64 `json:"priceOpen"`
		Shares       int64   `json:"shares"`
		Symbol       string  `json:"symbol"`
		Volume       int     `json:"volume"`
		VolumeAvg    int     `json:"volumeAvg"`
		Sector       string  `json:"sector"`
		SubSector    string  `json:"subSector"`
		Segment      string  `json:"segment"`
	}

	Client interface {
		Do(r *http.Request) (*http.Response, error)
	}

	Provider struct {
		Client Client
	}
)

func NewProvider(client Client) *Provider {
	return &Provider{
		Client: client,
	}
}

func (p Provider) LastInfo(ctx context.Context, symbol stock.Stock) (stock.Info, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", baseUrl, symbol), nil)
	if err != nil {
		return stock.Info{}, err
	}

	req = req.WithContext(ctx)
	res, err := p.Client.Do(req)
	if err != nil {
		return stock.Info{}, err
	}

	var data Response
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return stock.Info{}, err
	}

	return stock.Info{
		OpeningPrice: data.PriceOpen,
		MaxPrice:     data.High,
		MinPrice:     data.Low,
		LastPrice:    data.LastPrice,
		Change:       data.Change,
	}, nil
}
