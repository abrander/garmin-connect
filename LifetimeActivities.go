package connect

import (
	"errors"
)

// LifetimeActivities is describing a basic summary of all activities.
type LifetimeActivities struct {
	Activities    int     `json:"totalActivities"`    // The number of activities
	Distance      float64 `json:"totalDistance"`      // The total distance in meters
	Duration      float64 `json:"totalDuration"`      // The duration of all activities in seconds
	Calories      float64 `json:"totalCalories"`      // Energy in C
	ElevationGain float64 `json:"totalElevationGain"` // Total elevation gain in meters
}

// LifetimeActivities will return some aggregated data about all activities.
func (c *Client) LifetimeActivities(displayName string) (*LifetimeActivities, error) {
	URL := "https://connect.garmin.com/modern/proxy/userstats-service/statistics/" + displayName

	var proxy struct {
		Activities []LifetimeActivities `json:"userMetrics"`
	}

	err := c.getJSON(URL, &proxy)
	if err != nil {
		return nil, err
	}

	if len(proxy.Activities) != 1 {
		return nil, errors.New("unexpected data")
	}

	return &proxy.Activities[0], err
}
