package connect

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Activity describes a Garmin Connect activity.
type Activity struct {
	ID              int          `json:"activityId"`
	ActivityName    string       `json:"activityName"`
	Description     string       `json:"description"`
	StartLocal      Time         `json:"startTimeLocal"`
	StartGMT        Time         `json:"startTimeGMT"`
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

// RenameActivity can be used to rename an activity.
func (c *Client) RenameActivity(activityID int, newName string) error {
	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/activity-service/activity/%d", activityID)

	payload := struct {
		ID   int    `json:"activityId"`
		Name string `json:"activityName"`
	}{activityID, newName}

	body := bytes.NewBuffer(nil)
	enc := json.NewEncoder(body)
	err := enc.Encode(payload)
	if err != nil {
		return err
	}

	req, err := c.newRequest("PUT", URL, body)
	if err != nil {
		return err
	}

	req.Header.Add("nk", "NT")
	req.Header.Add("content-type", "application/json")

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode != 204 {
		return fmt.Errorf("HTTP call returned %d", resp.StatusCode)
	}

	return nil
}
