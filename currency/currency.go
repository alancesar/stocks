package currency

import (
	"fmt"
	"math"
	"strings"
)

type Currency struct {
	value float64
}

func NewFromFloat(value float64) Currency {
	return Currency{
		value: value,
	}
}

func (c Currency) Float64() float64 {
	return c.value
}

func (c Currency) String() string {
	raw := fmt.Sprintf("%.2f", math.Abs(c.value))
	raw = strings.ReplaceAll(raw, ".", ",")

	if c.value < 0 {
		return fmt.Sprintf("-(R$ %s)", raw)
	}
	return fmt.Sprintf("R$ %s", raw)
}
