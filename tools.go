package connect

import (
	"fmt"
	"time"
)

// date formats a time.Time as a date usable in the Garmin Connect API.
func formatDate(t time.Time) string {
	return fmt.Sprintf("%04d-%02d-%02d", t.Year(), t.Month(), t.Day())
}
