package operation

import (
	"context"
	"fmt"
	"io"
	"stocks/currency"
	"stocks/separator"
	"stocks/stock"
	"time"
)

const (
	Buy Type = iota
	Sell
)

type (
	Type int

	Repository interface {
		Create(ctx context.Context, operation Operation) error
		List(ctx context.Context) (List, error)
	}

	Operation struct {
		Symbol    stock.Symbol
		Type      Type
		Quantity  int
		UnitValue float64
		Date      time.Time
	}

	List []Operation
)

func (t Type) String() string {
	switch t {
	case Buy:
		return "BUY"
	case Sell:
		return "SELL"
	default:
		return ""
	}
}

func (o Operation) String() string {
	return fmt.Sprintf("Symbol=%-6s Type=%s Quantity=%d UnitValue=%s Date=%s",
		o.Symbol, o.Type, o.Quantity, currency.NewFromFloat(o.UnitValue), o.Date.Format("2006-01-02"))
}

func (l List) Print(writer io.Writer, sep separator.Separator) error {
	title := fmt.Sprintf("Symbol%sType%sQtd%sUnit. Value%sDate\n", sep, sep, sep, sep)
	if _, err := io.WriteString(writer, title); err != nil {
		return err
	}

	for _, operation := range l {
		line := printLine(operation, sep)

		if _, err := io.WriteString(writer, line); err != nil {
			return err
		}
	}

	return nil
}

func printLine(operation Operation, sep separator.Separator) string {
	switch sep {
	case separator.Tab:
		return printBeauty(operation, sep)
	default:
		return printRaw(operation, sep)
	}
}

func printRaw(operation Operation, sep separator.Separator) string {
	return fmt.Sprintf("%s%s%s%s%d%s%.2f%s%s\n",
		operation.Symbol, sep, operation.Type, sep, operation.Quantity, sep, operation.UnitValue, sep,
		operation.Date.Format("2006-01-02"))
}

func printBeauty(operation Operation, sep separator.Separator) string {
	return fmt.Sprintf("%s%s%s%s%d%s%s%s%s\n",
		operation.Symbol, sep, operation.Type, sep, operation.Quantity, sep, currency.NewFromFloat(operation.UnitValue),
		sep, operation.Date.Format("2006-01-02"))
}
