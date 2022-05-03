package main

import (
	"reflect"
	"stocks/date"
	"stocks/usecase"
	"testing"
	"time"
)

func TestCreateBuyRequest(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		args    args
		want    usecase.BuyRequest
		wantErr bool
	}{
		{
			name: "Should build request properly without data",
			args: args{
				args: []string{"STOCK", "10", "1.23"},
			},
			want: usecase.BuyRequest{
				Symbol:    "STOCK",
				Quantity:  10,
				UnitValue: 1.23,
				Date:      date.Trunc(time.Now()),
			},
			wantErr: false,
		},
		{
			name: "Should build request properly with data",
			args: args{
				args: []string{"STOCK", "10", "1.23", "2022-04-28"},
			},
			want: usecase.BuyRequest{
				Symbol:    "STOCK",
				Quantity:  10,
				UnitValue: 1.23,
				Date:      time.Date(2022, 4, 28, 0, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "Should return error if is missing args",
			args: args{
				args: []string{"STOCK"},
			},
			want:    usecase.BuyRequest{},
			wantErr: true,
		},
		{
			name: "Should return error if quantity is invalid",
			args: args{
				args: []string{"STOCK", "abc", "1.23", "2022-04-28"},
			},
			want:    usecase.BuyRequest{},
			wantErr: true,
		},
		{
			name: "Should return error if value is invalid",
			args: args{
				args: []string{"STOCK", "10", "abc", "2022-04-28"},
			},
			want:    usecase.BuyRequest{},
			wantErr: true,
		},
		{
			name: "Should return error if date is invalid",
			args: args{
				args: []string{"STOCK", "10", "1.23", "abc"},
			},
			want:    usecase.BuyRequest{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateBuyRequest(tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateBuyRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateBuyRequest() got = %v, want %v", got, tt.want)
			}
		})
	}
}
