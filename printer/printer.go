package printer

import (
	"io"
	"stocks/separator"
)

type Printer interface {
	Print(writer io.Writer, sep separator.Separator) error
}
