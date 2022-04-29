package printer

import (
	"fmt"
	"io"
	"stocks/operation"
)

func Print(writer io.Writer, report operation.Report) error {
	if _, err := io.WriteString(writer, "Stock\tQtd.\tAvg. Price\tLast Price\tGain/Loss\n"); err != nil {
		return err
	}

	for _, e := range report.Summary {
		line := fmt.Sprintf("%s\t%d\t%s\t%s\t%s\n",
			e.Stock, e.Quantity, e.AveragePrice, e.LastPrice, e.GainLoss())

		if _, err := io.WriteString(writer, line); err != nil {
			return err
		}
	}

	return nil
}
