package date

import "time"

const (
	layout = "2006-01-02"
)

func Parse(raw string) (time.Time, error) {
	switch raw {
	case "today":
		return Trunc(time.Now()), nil
	case "yesterday":
		return Trunc(time.Now().Add(-24 * time.Hour)), nil
	default:
		return time.Parse(layout, raw)
	}
}
func Trunc(source time.Time) time.Time {
	return time.Date(source.Year(), source.Month(), source.Day(), 0, 0, 0, 0, time.UTC)
}
