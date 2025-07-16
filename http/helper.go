package http

import (
	"fmt"
	"math"

	"github.com/dustin/go-humanize"
)

func FormatTZS(v float64) string {
	// We don’t expect decimals; round to nearest shilling
	n := int64(math.Round(v))
	return fmt.Sprintf("%s TZS", humanize.Comma(n)) // → “35,000 TZS”
}
