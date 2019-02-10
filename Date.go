package connect

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Date represents a single day in Garmin Connect.
type Date struct {
	Year       int
	Month      int
	DayOfMonth int
}

// Time returns a time.Time for usage in other packages.
func (d Date) Time() time.Time {
	return time.Date(d.Year, time.Month(d.Month), d.DayOfMonth, 0, 0, 0, 0, time.UTC)
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *Date) UnmarshalJSON(value []byte) error {
	// Sometimes dates are transferred as milliseconds since epoch :-/
	i, err := strconv.ParseInt(string(value), 10, 64)
	if err == nil {
		t := time.Unix(i/1000, 0)

		d.Year = t.Year()
		d.Month = int(t.Month())
		d.DayOfMonth = t.Day()

		return nil
	}

	var blip string
	err = json.Unmarshal(value, &blip)
	if err != nil {
		return err
	}

	_, err = fmt.Sscanf(blip, "%04d-%02d-%02d", &d.Year, &d.Month, &d.DayOfMonth)
	if err != nil {
		return err
	}

	return nil
}

// String implements Stringer.
func (d *Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.DayOfMonth)
}
