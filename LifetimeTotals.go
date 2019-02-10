package connect

// LifetimeTotals is ligetime statistics for the Connect user.
type LifetimeTotals struct {
	ProfileID      int     `json:"userProfileId"`
	ActiveDays     int     `json:"totalActiveDays"`
	Calories       float64 `json:"totalCalories"`
	Distance       int     `json:"totalDistance"`
	GoalsMetInDays int     `json:"totalGoalsMetInDays"`
	Steps          int     `json:"totalSteps"`
}

// LifetimeTotals returns some lifetime statistics for displayName.
func (c *Client) LifetimeTotals(displayName string) (*LifetimeTotals, error) {
	URL := "https://connect.garmin.com/modern/proxy/usersummary-service/stats/connectLifetimeTotals/" + displayName

	totals := new(LifetimeTotals)

	err := c.getJSON(URL, totals)
	if err != nil {
		return nil, err
	}

	return totals, err
}
