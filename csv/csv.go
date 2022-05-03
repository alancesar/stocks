package csv

import (
	"encoding/csv"
	"io"
)

func Import[T any](reader io.Reader, fn func([]string) (T, error)) ([]T, error) {
	r := csv.NewReader(reader)

	var output []T
	for {
		record, err := r.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		if item, err := fn(record); err != nil {
			return nil, err
		} else {
			output = append(output, item)
		}
	}

	return output, nil
}
