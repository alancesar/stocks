package currency

import (
	"reflect"
	"testing"
)

func TestCurrency_Float64(t *testing.T) {
	type fields struct {
		value float64
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name: "Should return float value",
			fields: fields{
				value: 1.23,
			},
			want: 1.23,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Currency{
				value: tt.fields.value,
			}
			if got := c.Float64(); got != tt.want {
				t.Errorf("Float64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCurrency_String(t *testing.T) {
	type fields struct {
		value float64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Should print positive values properly",
			fields: fields{
				value: -1.23,
			},
			want: "-(R$ 1,23)",
		},
		{
			name: "Should print negative values properly",
			fields: fields{
				value: 1.23,
			},
			want: "R$ 1,23",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Currency{
				value: tt.fields.value,
			}
			if got := c.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewFromFloat(t *testing.T) {
	type args struct {
		value float64
	}
	tests := []struct {
		name string
		args args
		want Currency
	}{
		{
			name: "Should create from float properly",
			args: args{
				value: 1.23,
			},
			want: Currency{
				value: 1.23,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFromFloat(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFromFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}
