package connect

import (
	"fmt"
	"time"
)

// ActivityHrZones describes the heart-rate zones during an activity.
type ActivityHrZones struct {
	TimeInZone      time.Duration `json:"secsInZone"`
	ZoneLowBoundary int           `json:"zoneLowBoundary"`
	ZoneNumber      int           `json:"zoneNumber"`
}

// ActivityHrZones returns the reported heart-rate zones for an activity.
func (c *Client) ActivityHrZones(activityID int) ([]ActivityHrZones, error) {
	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/activity-service/activity/%d/hrTimeInZones",
		activityID,
	)

	var proxy []struct {
		TimeInZone      float64 `json:"secsInZone"`
		ZoneLowBoundary int     `json:"zoneLowBoundary"`
		ZoneNumber      int     `json:"zoneNumber"`
	}

	err := c.getJSON(URL, &proxy)
	if err != nil {
		return nil, err
	}

	zones := make([]ActivityHrZones, len(proxy))

	for i, p := range proxy {
		zones[i].TimeInZone = time.Duration(p.TimeInZone * float64(time.Second))
		zones[i].ZoneLowBoundary = p.ZoneLowBoundary
		zones[i].ZoneNumber = p.ZoneNumber
	}

	return zones, nil
}
