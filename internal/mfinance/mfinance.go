package mfinance

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"stocks/stock"
)

const (
	baseUrl = "https://mfinance.com.br/api/v1"
)

type (
	Info struct {
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

	Details struct {
		Symbol    string `json:"symbol"`
		Type      string `json:"type"`
		SubSector string `json:"subSector"`
		Segment   string `json:"segment"`
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

func (p Provider) LastInfo(ctx context.Context, symbol stock.Symbol) (stock.Info, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/stocks/%s", baseUrl, symbol), nil)
	if err != nil {
		return stock.Info{}, err
	}

	req = req.WithContext(ctx)
	res, err := p.Client.Do(req)
	if err != nil {
		return stock.Info{}, err
	}

	info, err := decode[Info](res)
	if err != nil {
		return stock.Info{}, err
	}

	return stock.Info{
		Symbol:       symbol,
		OpeningPrice: info.PriceOpen,
		MaxPrice:     info.High,
		MinPrice:     info.Low,
		LastPrice:    info.LastPrice,
		Change:       info.Change,
	}, nil
}

func (p Provider) Details(ctx context.Context, symbol stock.Symbol) (stock.Details, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/stocks/details/%s", baseUrl, symbol), nil)
	if err != nil {
		return stock.Details{}, err
	}

	req = req.WithContext(ctx)
	res, err := p.Client.Do(req)
	if err != nil {
		return stock.Details{}, err
	}

	details, err := decode[Details](res)
	if err != nil {
		return stock.Details{}, err
	} else if details.Symbol == "" {
		return stock.Details{}, fmt.Errorf("%s not found", symbol)
	}

	return stock.Details{
		Symbol:  symbol,
		Type:    details.Type,
		Sector:  details.SubSector,
		Segment: details.Segment,
	}, nil
}

func decode[T any](response *http.Response) (T, error) {
	defer func() {
		_ = response.Body.Close()
	}()

	var output T

	if response.StatusCode != http.StatusOK {
		return output, fmt.Errorf("invalid status code: %d", response.StatusCode)
	}

	if err := json.NewDecoder(response.Body).Decode(&output); err != nil {
		return output, err
	}

	return output, nil
}
