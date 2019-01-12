package connect

import (
	"bytes"
	"encoding/json"
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
	SourceType        string  `json:"sourceType"`
}

// WeightAverage is aggregated weight data for a specific period.
type WeightAverage struct {
	Weightin
	From  int `json:"from"`
	Until int `json:"until"`
}

// LatestWeight will retrieve the latest weight by date.
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

// Weightins will retrieve all weight ins between startDate and endDate. A
// summary is provided as well. This summary is calculated by Garmin Connect.
func (c *Client) Weightins(startDate time.Time, endDate time.Time) (*WeightAverage, []Weightin, error) {
	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/weight-service/weight/dateRange?startDate=%s&endDate=%s",
		formatDate(startDate),
		formatDate(endDate))

	var proxy struct {
		DateWeightList []Weightin     `json:"dateWeightList"`
		TotalAverage   *WeightAverage `json:"totalAverage"`
	}

	err := c.getJSON(URL, &proxy)
	if err != nil {
		return nil, nil, err
	}

	return proxy.TotalAverage, proxy.DateWeightList, nil
}

// DeleteWeightin will delete all biometric data for date.
func (c *Client) DeleteWeightin(date time.Time) error {
	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/biometric-service/biometric/%s", formatDate(date))

	req, err := c.newRequest("DELETE", URL, nil)
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}

	resp.Body.Close()

	return nil
}

// AddUserWeight will add a manual weight in. weight is in grams to match
// Weightin.
func (c *Client) AddUserWeight(date time.Time, weight float64) error {
	URL := "https://connect.garmin.com/modern/proxy/weight-service/user-weight"
	payload := struct {
		Date    string  `json:"date"`
		UnitKey string  `json:"unitKey"`
		Value   float64 `json:"value"`
	}{
		Date:    formatDate(date),
		UnitKey: "kg",
		Value:   weight / 1000.0,
	}

	body := bytes.NewBuffer(nil)
	enc := json.NewEncoder(body)
	err := enc.Encode(payload)
	if err != nil {
		return err
	}

	req, err := c.newRequest("POST", URL, body)
	if err != nil {
		return err
	}
	req.Header.Add("nk", "NT") // Yep. This is needed. No idea what it does.
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
