package connect

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// Activity describes a Garmin Connect activity.
type Activity struct {
	ID               int          `json:"activityId"`
	ActivityName     string       `json:"activityName"`
	Description      string       `json:"description"`
	StartLocal       Time         `json:"startTimeLocal"`
	StartGMT         Time         `json:"startTimeGMT"`
	ActivityType     ActivityType `json:"activityType"`
	Distance         float64      `json:"distance"` // meter
	Duration         float64      `json:"duration"`
	ElapsedDuration  float64      `json:"elapsedDuration"`
	MovingDuration   float64      `json:"movingDuration"`
	AverageSpeed     float64      `json:"averageSpeed"`
	MaxSpeed         float64      `json:"maxSpeed"`
	OwnerID          int          `json:"ownerId"`
	Calories         float64      `json:"calories"`
	AverageHeartRate float64      `json:"averageHR"`
	MaxHeartRate     float64      `json:"maxHR"`
	DeviceID         int          `json:"deviceId"`
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

const (
	// FormatFIT is the "original" Garmin format. Please note that this will be written as a ZIP file (!).
	FormatFIT = iota

	// FormatTCX is Training Center XML (TCX) format.
	FormatTCX

	// FormatGPX will export as GPX - the GPS Exchange Format.
	FormatGPX

	// FormatKML will export KML files compatible with Google Earth.
	FormatKML

	// FormatCSV will export splits as CSV.
	FormatCSV

	formatMax
)

// ExportActivity will export an activity from Connect. The activity will be written til w.
func (c *Client) ExportActivity(id int, w io.Writer, format int) error {
	formatTable := [formatMax]string{
		"https://connect.garmin.com/modern/proxy/download-service/files/activity/%d",
		"https://connect.garmin.com/modern/proxy/download-service/export/tcx/activity/%d",
		"https://connect.garmin.com/modern/proxy/download-service/export/gpx/activity/%d",
		"https://connect.garmin.com/modern/proxy/download-service/export/kml/activity/%d",
		"https://connect.garmin.com/modern/proxy/download-service/export/csv/activity/%d",
	}

	if format >= formatMax || format < FormatFIT {
		return errors.New("Invalid format")
	}

	URL := fmt.Sprintf(formatTable[format], id)

	return c.download(URL, w)
}
