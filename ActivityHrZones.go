package connect

import (
	"fmt"
	"time"
)

// ActivityHrZones describes the heart-rate zones during an activity.
type ActivityHrZones struct {
	SecsInZone      time.Duration `json:"secsInZone"`
	ZoneLowBoundary int           `json:"zoneLowBoundary"`
	ZoneNumber      int           `json:"zoneNumber"`
}

// ActivityHrZones returns the reported heart-rate zones for an activity.
func (c *Client) ActivityHrZones(activityID int) ([]ActivityHrZones, error) {
	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/activity-service/activity/%d/hrTimeInZones",
		activityID,
	)

	var hrZones []ActivityHrZones

	err := c.getJSON(URL, &hrZones)
	if err != nil {
		return nil, err
	}

	return hrZones, nil
}
