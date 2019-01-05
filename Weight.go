package connect

import (
	"fmt"
	"time"
)

// Weightin is a single weight event.
type Weightin struct {
	Date              int     `json:"date"`
	Version           int     `json:"version"`
	Weight            float64 `json:"weight"`     // gram
	BMI               float64 `json:"bmi"`        // weight / height²
	BodyFatPercentage float64 `json:"bodyFat"`    // percent
	BodyWater         float64 `json:"bodyWater"`  // kilogram
	BoneMass          int     `json:"boneMass"`   // gram
	MuscleMass        int     `json:"muscleMass"` // gram
	//"physiqueRating": null,
	//"visceralFat": null,
	//"metabolicAge": null,
	//"caloricIntake": null,
	SourceType string `json:"sourceType"`
}

// LatestWeight will retrieve the latest weight in by date.
func (c *Client) LatestWeight(date time.Time) (*Weightin, error) {
	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/weight-service/weight/latest?date=%04d-%02d-%02d",
		date.Year(),
		date.Month(),
		date.Day())

	wi := new(Weightin)

	err := c.getJSON(URL, wi)
	if err != nil {
		return nil, err
	}

	return wi, nil
}
