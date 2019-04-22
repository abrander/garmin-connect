package connect

import (
	"archive/zip"
	"bytes"
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

// Activity will retrieve details about an activity.
func (c *Client) Activity(activityID int) (*Activity, error) {
	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/activity-service/activity/%d",
		activityID,
	)

	activity := new(Activity)

	err := c.getJSON(URL, &activity)
	if err != nil {
		return nil, err
	}

	return activity, nil
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

	return c.write("PUT", URL, payload, 204)
}

const (
	// FormatFIT is the "original" Garmin format.
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
		return errors.New("invalid format")
	}

	URL := fmt.Sprintf(formatTable[format], id)

	// To unzip FIT files on-the-fly, we treat them specially.
	if format == FormatFIT {
		buffer := bytes.NewBuffer(nil)

		err := c.download(URL, buffer)
		if err != nil {
			return err
		}

		z, err := zip.NewReader(bytes.NewReader(buffer.Bytes()), int64(buffer.Len()))
		if err != nil {
			return err
		}

		if len(z.File) != 1 {
			return fmt.Errorf("%d files found in FIT archive, 1 expected", len(z.File))
		}

		src, err := z.File[0].Open()
		if err != nil {
			return err
		}
		defer src.Close()

		_, err = io.Copy(w, src)
		return err
	}

	return c.download(URL, w)
}

// DeleteActivity will permanently delete an activity.
func (c *Client) DeleteActivity(id int) error {
	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/activity-service/activity/%d", id)

	return c.write("DELETE", URL, nil, 0)
}
