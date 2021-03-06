package asset

import (
	"bytes"
	"reflect"
	"stocks/currency"
	"stocks/separator"
	"stocks/stock"
	"testing"
)

func TestAsset_Balance(t *testing.T) {
	type fields struct {
		Investment currency.Currency
		Settled    currency.Currency
	}
	tests := []struct {
		name   string
		fields fields
		want   currency.Currency
	}{
		{
			name: "Should calculate negative balance properly",
			fields: fields{
				Investment: currency.NewFromFloat(100),
				Settled:    currency.NewFromFloat(0),
			},
			want: currency.NewFromFloat(-100),
		},
		{
			name: "Should calculate zero balance properly",
			fields: fields{
				Investment: currency.NewFromFloat(100),
				Settled:    currency.NewFromFloat(200),
			},
			want: currency.NewFromFloat(100),
		},
		{
			name: "Should calculate positive balance properly",
			fields: fields{
				Investment: currency.NewFromFloat(100),
				Settled:    currency.NewFromFloat(100),
			},
			want: currency.NewFromFloat(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Asset{
				Investment: tt.fields.Investment,
				Settled:    tt.fields.Settled,
			}
			if got := s.Balance(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Balance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAsset_GainLoss(t *testing.T) {
	type fields struct {
		Symbol       stock.Symbol
		Quantity     int
		AveragePrice currency.Currency
		LastPrice    currency.Currency
		Investment   currency.Currency
		Settled      currency.Currency
	}
	tests := []struct {
		name   string
		fields fields
		want   currency.Currency
	}{
		{
			name: "Should calculate gain for stocks valuation using current value",
			fields: fields{
				Quantity:   10,
				LastPrice:  currency.NewFromFloat(12),
				Investment: currency.NewFromFloat(100),
			},
			want: currency.NewFromFloat(20),
		},
		{
			name: "Should calculate gain for stocks in sell operations",
			fields: fields{
				Investment: currency.NewFromFloat(100),
				Settled:    currency.NewFromFloat(120),
			},
			want: currency.NewFromFloat(20),
		},
		{
			name: "Should calculate gain for stocks using current valuation and sell operations",
			fields: fields{
				Quantity:   5,
				LastPrice:  currency.NewFromFloat(12),
				Investment: currency.NewFromFloat(100),
				Settled:    currency.NewFromFloat(60),
			},
			want: currency.NewFromFloat(20),
		},
		{
			name: "Should calculate loss for stocks devaluation",
			fields: fields{
				Quantity:   10,
				LastPrice:  currency.NewFromFloat(8),
				Investment: currency.NewFromFloat(100),
			},
			want: currency.NewFromFloat(-20),
		},
		{
			name: "Should calculate loss for stocks in sell operations",
			fields: fields{
				Investment: currency.NewFromFloat(100),
				Settled:    currency.NewFromFloat(80),
			},
			want: currency.NewFromFloat(-20),
		},
		{
			name: "Should calculate loss for stocks using current valuation and sell operations",
			fields: fields{
				Quantity:   5,
				LastPrice:  currency.NewFromFloat(8),
				Investment: currency.NewFromFloat(100),
				Settled:    currency.NewFromFloat(40),
			},
			want: currency.NewFromFloat(-20),
		},
		{
			name: "Should calculate loss for stocks using current valuation and sell operations",
			fields: fields{
				Quantity:   5,
				LastPrice:  currency.NewFromFloat(8),
				Investment: currency.NewFromFloat(120),
				Settled:    currency.NewFromFloat(80),
			},
			want: currency.NewFromFloat(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Asset{
				Quantity:   tt.fields.Quantity,
				LastPrice:  tt.fields.LastPrice,
				Investment: tt.fields.Investment,
				Settled:    tt.fields.Settled,
			}
			if got := s.GainLoss(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GainLoss() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAssets_Balance(t *testing.T) {
	type fields struct {
		assets Assets
	}
	tests := []struct {
		name   string
		fields fields
		want   currency.Currency
	}{
		{
			name: "Should calculate assets report balance properly",
			fields: fields{
				assets: Assets{
					{
						Investment: currency.NewFromFloat(200),
					},
					{
						Investment: currency.NewFromFloat(100),
						Settled:    currency.NewFromFloat(50),
					},
					{
						Investment: currency.NewFromFloat(100),
						Settled:    currency.NewFromFloat(200),
					},
				},
			},
			want: currency.NewFromFloat(-150),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assets := tt.fields.assets
			if got := assets.Balance(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Balance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAssets_GainLoss(t *testing.T) {
	type fields struct {
		assets Assets
	}
	tests := []struct {
		name   string
		fields fields
		want   currency.Currency
	}{
		{
			name: "Should calculate assets gain / loss properly",
			fields: fields{
				assets: Assets{
					{
						Quantity:   10,
						LastPrice:  currency.NewFromFloat(12),
						Investment: currency.NewFromFloat(100),
					},
					{
						Investment: currency.NewFromFloat(100),
						Settled:    currency.NewFromFloat(120),
					},
					{
						Quantity:   5,
						LastPrice:  currency.NewFromFloat(12),
						Investment: currency.NewFromFloat(100),
						Settled:    currency.NewFromFloat(60),
					},
				},
			},
			want: currency.NewFromFloat(60),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assets := tt.fields.assets
			if got := assets.GainLoss(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GainLoss() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAssets_Print(t *testing.T) {
	type fields struct {
		assets Assets
	}
	type args struct {
		sep separator.Separator
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantWriter string
		wantErr    bool
	}{
		{
			name: "Should print assets report properly with comma",
			fields: fields{
				assets: Assets{
					{
						Symbol:       "STOCK1",
						Quantity:     8,
						AveragePrice: currency.NewFromFloat(10),
						LastPrice:    currency.NewFromFloat(12),
						Investment:   currency.NewFromFloat(100),
						Settled:      currency.NewFromFloat(0),
					},
					{
						Symbol:       "STOCK2",
						Quantity:     6,
						AveragePrice: currency.NewFromFloat(9.8),
						LastPrice:    currency.NewFromFloat(21),
						Investment:   currency.NewFromFloat(220),
						Settled:      currency.NewFromFloat(160),
					},
				},
			},
			args: args{
				sep: separator.Comma,
			},
			wantWriter: "Symbol,Qtd.,Avg. Price,Last Price,Gain/Loss\nSTOCK1,8,R$ 10,00,R$ 12,00,-(R$ 4,00)\nSTOCK2,6,R$ 9,80,R$ 21,00,R$ 66,00\n",
			wantErr:    false,
		},
		{
			name: "Should print assets report properly with tab",
			fields: fields{
				assets: Assets{
					{
						Symbol:       "STOCK1",
						Quantity:     8,
						AveragePrice: currency.NewFromFloat(10),
						LastPrice:    currency.NewFromFloat(12),
						Investment:   currency.NewFromFloat(100),
						Settled:      currency.NewFromFloat(0),
					},
					{
						Symbol:       "STOCK2",
						Quantity:     6,
						AveragePrice: currency.NewFromFloat(9.8),
						LastPrice:    currency.NewFromFloat(21),
						Investment:   currency.NewFromFloat(220),
						Settled:      currency.NewFromFloat(160),
					},
				},
			},
			args: args{
				sep: separator.Tab,
			},
			wantWriter: "Symbol\tQtd.\tAvg. Price\tLast Price\tGain/Loss\nSTOCK1\t8\tR$ 10,00\tR$ 12,00\t-(R$ 4,00)\nSTOCK2\t6\tR$ 9,80\tR$ 21,00\tR$ 66,00\n",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assets := tt.fields.assets
			writer := &bytes.Buffer{}
			err := assets.Print(writer, tt.args.sep)
			if (err != nil) != tt.wantErr {
				t.Errorf("Print() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("Print() gotWriter = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}
