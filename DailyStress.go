package connect

import (
	"fmt"
	"time"
)

// StressPoint is a measured stress level at a point in time.
type StressPoint struct {
	Timestamp time.Time
	Value     int
}

// DailyStress is a stress reading for a single day.
type DailyStress struct {
	UserProfilePK int    `json:"userProfilePK"`
	CalendarDate  string `json:"calendarDate"`
	StartGMT      Time   `json:"startTimestampGMT"`
	EndGMT        Time   `json:"endTimestampGMT"`
	StartLocal    Time   `json:"startTimestampLocal"`
	EndLocal      Time   `json:"endTimestampLocal"`
	Max           int    `json:"maxStressLevel"`
	Average       int    `json:"avgStressLevel"`
	Values        []StressPoint
}

// DailyStress will retrieve stress levels for date.
func (c *Client) DailyStress(date time.Time) (*DailyStress, error) {
	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/wellness-service/wellness/dailyStress/%s",
		formatDate(date))

	if !c.authenticated() {
		return nil, ErrNotAuthenticated
	}

	// We use a proxy object to deserialize the values to proper Go types.
	var proxy struct {
		DailyStress
		StressValuesArray [][2]int64 `json:"stressValuesArray"`
	}

	err := c.getJSON(URL, &proxy)
	if err != nil {
		return nil, err
	}

	ret := &proxy.DailyStress
	ret.Values = make([]StressPoint, len(proxy.StressValuesArray))

	for i, point := range proxy.StressValuesArray {
		ret.Values[i].Timestamp = time.Unix(point[0]/1000, 0)
		ret.Values[i].Value = int(point[1])
	}

	return &proxy.DailyStress, nil
}
