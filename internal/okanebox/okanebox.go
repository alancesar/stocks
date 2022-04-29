package okanebox

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"stocks/stock"
	"time"
)

const (
	baseUrl = "https://www.okanebox.com.br/api/acoes/ultima"
)

type (
	Response struct {
		ClosingDate  string  `json:"DATPRG"`
		OpeningPrice float64 `json:"PREABE"`
		MaxPrice     float64 `json:"PREMAX"`
		MinPrice     float64 `json:"PREMIN"`
		MediumPrice  float64 `json:"PREMED"`
		LastPrice    float64 `json:"PREULT"`
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

	parsedDate, err := time.Parse("2006-01-02T15:04:05", data.ClosingDate)
	if err != nil {
		return stock.Info{}, err
	}

	return stock.Info{
		ClosingDate:  parsedDate,
		OpeningPrice: data.OpeningPrice,
		MaxPrice:     data.MaxPrice,
		MinPrice:     data.MinPrice,
		MediumPrice:  data.MediumPrice,
		LastPrice:    data.LastPrice,
	}, nil
}
