package connect

import (
	"fmt"
	"time"
)

// Weightin is a single weight event.
type Weightin struct {
	Date              Date    `json:"date"`
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

	// An alternative endpoint for weight info this can be found here:
	// https://connect.garmin.com/modern/proxy/userprofile-service/userprofile/personal-information/weightWithOutbound?from=1556359100000&until=1556611800000

	if !c.authenticated() {
		return nil, nil, ErrNotAuthenticated
	}

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

	if !c.authenticated() {
		return ErrNotAuthenticated
	}

	return c.write("DELETE", URL, nil, 204)
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

	return c.write("POST", URL, payload, 204)
}

// WeightByDate retrieves the weight of date if available. If no weight data
// for date exists, it will return ErrNotFound.
func (c *Client) WeightByDate(date time.Time) (Time, float64, error) {
	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/biometric-service/biometric/weightByDate?date=%s",
		formatDate(date))

	if !c.authenticated() {
		return Time{}, 0.0, ErrNotAuthenticated
	}

	var proxy []struct {
		TimeStamp Time    `json:"weightDate"`
		Weight    float64 `json:"weight"` // gram
	}

	err := c.getJSON(URL, &proxy)
	if err != nil {
		return Time{}, 0.0, err
	}

	if len(proxy) < 1 {
		return Time{}, 0.0, ErrNotFound
	}

	return proxy[0].TimeStamp, proxy[0].Weight, nil
}

// WeightGoal will list the users weight goal if any. If displayName is empty,
// the currently authenticated user will be used.
func (c *Client) WeightGoal(displayName string) (*Goal, error) {
	goals, err := c.Goals(displayName, 4)
	if err != nil {
		return nil, err
	}

	if len(goals) < 1 {
		return nil, ErrNotFound
	}

	return &goals[0], nil
}

// SetWeightGoal will set a new weight goal.
func (c *Client) SetWeightGoal(goal int) error {
	if !c.authenticated() || c.Profile == nil {
		return ErrNotAuthenticated
	}

	g := Goal{
		Created:   Today(),
		Start:     Today(),
		GoalType:  4,
		ProfileID: c.Profile.ProfileID,
		Value:     goal,
	}

	goals, err := c.Goals("", 4)
	if err != nil {
		return err
	}

	if len(goals) >= 1 {
		g.ID = goals[0].ID
		return c.UpdateGoal("", g)
	}

	return c.AddGoal(c.Profile.DisplayName, g)
}
