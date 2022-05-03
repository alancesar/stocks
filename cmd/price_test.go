package main

import (
	"reflect"
	"stocks/stock"
	"testing"
)

func TestCreatePriceRequest(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		args    args
		want    stock.Symbol
		wantErr bool
	}{
		{
			name: "Should create request properly",
			args: args{
				args: []string{"STOCK"},
			},
			want:    stock.Symbol("STOCK"),
			wantErr: false,
		},
		{
			name: "Should return error if symbol arg is missing",
			args: args{
				args: []string{""},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreatePriceRequest(tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreatePriceRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreatePriceRequest() got = %v, want %v", got, tt.want)
			}
		})
	}
}
