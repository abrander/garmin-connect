package connect

import (
	"time"
)

// Timezone represents a timezone in Garmin Connect.
type Timezone struct {
	ID        int     `json:"unitId"`
	Key       string  `json:"unitKey"`
	GMTOffset float64 `json:"gmtOffset"`
	DSTOffset float64 `json:"dstOffset"`
	Group     int     `json:"groupNumber"`
	TimeZone  string  `json:"timeZone"`
}

// Location will (try to) return a location for use with time.Time functions.
func (t *Timezone) Location() (*time.Location, error) {
	return time.LoadLocation(t.Key)
}
