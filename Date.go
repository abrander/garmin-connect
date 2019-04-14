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
	Month      time.Month
	DayOfMonth int
}

// Time returns a time.Time for usage in other packages.
func (d Date) Time() time.Time {
	return time.Date(d.Year, d.Month, d.DayOfMonth, 0, 0, 0, 0, time.UTC)
}

// UnmarshalJSON implements json.Unmarshaler.
func (d *Date) UnmarshalJSON(value []byte) error {
	if string(value) == "null" {
		return nil
	}

	// Sometimes dates are transferred as milliseconds since epoch :-/
	i, err := strconv.ParseInt(string(value), 10, 64)
	if err == nil {
		t := time.Unix(i/1000, 0)

		d.Year, d.Month, d.DayOfMonth = t.Date()

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

// MarshalJSON implements json.Marshaler.
func (d Date) MarshalJSON() ([]byte, error) {
	// To better support the Garmin API we marshal the empty value as null.
	if d.Year == 0 && d.Month == 0 && d.DayOfMonth == 0 {
		return []byte("null"), nil
	}

	return []byte(fmt.Sprintf("\"%04d-%02d-%02d\"", d.Year, d.Month, d.DayOfMonth)), nil
}

// ParseDate will parse a date in the format yyyy-mm-dd.
func ParseDate(in string) (Date, error) {
	d := Date{}

	_, err := fmt.Sscanf(in, "%04d-%02d-%02d", &d.Year, &d.Month, &d.DayOfMonth)

	return d, err
}

// String implements Stringer.
func (d Date) String() string {
	if d.Year == 0 && d.Month == 0 && d.DayOfMonth == 0 {
		return "-"
	}

	return fmt.Sprintf("%04d-%02d-%02d", d.Year, d.Month, d.DayOfMonth)
}

// Today will return a Date set to today.
func Today() Date {
	d := Date{}

	d.Year, d.Month, d.DayOfMonth = time.Now().Date()

	return d
}
