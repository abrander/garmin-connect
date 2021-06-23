package connect

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
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

// ExportActivity will export an activity from Connect. The activity will be written til w.
func (c *Client) ExportActivity(id int, w io.Writer, format ActivityFormat) error {
	formatTable := [activityFormatMax]string{
		"https://connect.garmin.com/modern/proxy/download-service/files/activity/%d",
		"https://connect.garmin.com/modern/proxy/download-service/export/tcx/activity/%d",
		"https://connect.garmin.com/modern/proxy/download-service/export/gpx/activity/%d",
		"https://connect.garmin.com/modern/proxy/download-service/export/kml/activity/%d",
		"https://connect.garmin.com/modern/proxy/download-service/export/csv/activity/%d",
	}

	if format >= activityFormatMax || format < ActivityFormatFIT {
		return errors.New("invalid format")
	}

	URL := fmt.Sprintf(formatTable[format], id)

	// To unzip FIT files on-the-fly, we treat them specially.
	if format == ActivityFormatFIT {
		buffer := bytes.NewBuffer(nil)

		err := c.Download(URL, buffer)
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

	return c.Download(URL, w)
}

// ImportActivity will import an activity into Garmin Connect. The activity
// will be read from file.
func (c *Client) ImportActivity(file io.Reader, format ActivityFormat) (int, error) {
	URL := "https://connect.garmin.com/modern/proxy/upload-service/upload/." + format.Extension()

	switch format {
	case ActivityFormatFIT, ActivityFormatTCX, ActivityFormatGPX:
		// These are ok.
	default:
		return 0, fmt.Errorf("%s is not supported for import", format.Extension())
	}

	formData := bytes.Buffer{}
	writer := multipart.NewWriter(&formData)
	defer writer.Close()

	activity, err := writer.CreateFormFile("file", "activity."+format.Extension())
	if err != nil {
		return 0, err
	}

	_, err = io.Copy(activity, file)
	if err != nil {
		return 0, err
	}

	writer.Close()

	req, err := c.newRequest("POST", URL, &formData)
	if err != nil {
		return 0, err
	}

	req.Header.Add("content-type", writer.FormDataContentType())

	resp, err := c.do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Implement enough of the response to satisfy our needs.
	var response struct {
		ImportResult struct {
			Successes []struct {
				InternalID int `json:"internalId"`
			} `json:"successes"`

			Failures []struct {
				Messages []struct {
					Content string `json:"content"`
				} `json:"messages"`
			} `json:"failures"`
		} `json:"detailedImportResult"`
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return 0, err
	}

	// This is ugly.
	if len(response.ImportResult.Failures) > 0 {
		messages := make([]string, 0, 10)
		for _, f := range response.ImportResult.Failures {
			for _, m := range f.Messages {
				messages = append(messages, m.Content)
			}
		}

		return 0, errors.New(strings.Join(messages, "; "))
	}

	if resp.StatusCode != 201 {
		return 0, fmt.Errorf("%d: %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	if len(response.ImportResult.Successes) != 1 {
		return 0, Error("cannot parse response, no failures and no successes..?")
	}

	return response.ImportResult.Successes[0].InternalID, nil
}

// DeleteActivity will permanently delete an activity.
func (c *Client) DeleteActivity(id int) error {
	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/activity-service/activity/%d", id)

	return c.write("DELETE", URL, nil, 0)
}
