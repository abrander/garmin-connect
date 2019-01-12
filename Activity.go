package connect

import (
	"fmt"
)

// Activity describes a Garmin Connect activity.
type Activity struct {
	ActivityID      int          `json:"activityId"`
	ActivityName    string       `json:"activityName"`
	Description     string       `json:"description"`
	StartTimeLocal  string       `json:"startTimeLocal"`
	StartTimeGMT    string       `json:"startTimeGMT"`
	ActivityType    ActivityType `json:"activityType"`
	Distance        float64      `json:"distance"` // meter
	Duration        float64      `json:"duration"`
	ElapsedDuration float64      `json:"elapsedDuration"`
	MovingDuration  float64      `json:"movingDuration"`
	AverageSpeed    float64      `json:"averageSpeed"`
	MaxSpeed        float64      `json:"maxSpeed"`
}

// ActivityType describes the type of activity.
type ActivityType struct {
	TypeID       int    `json:"typeId"`
	TypeKey      string `json:"typeKey"`
	ParentTypeID int    `json:"parentTypeId"`
	SortOrder    int    `json:"sortOrder"`
}

// Activities will list activities for displayName. If displayName is empty,
// the authenticated user will be used.
func (c *Client) Activities(displayName string, start int, limit int) ([]Activity, error) {
	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/activitylist-service/activities/%s?start=%d&limit=%d", displayName, start, limit)

	if !c.authenticated() && displayName == "" {
		return nil, ErrNotAuthenticated
	}

	var proxy struct {
		List []Activity `json:"activityList"`
	}

	err := c.getJSON(URL, &proxy)
	if err != nil {
		return nil, err
	}

	return proxy.List, nil
}
