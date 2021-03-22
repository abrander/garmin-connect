package connect

import (
	"fmt"
)

// Goal represents a fitness or health goal.
type Goal struct {
	ID           int64    `json:"id"`
	ProfileID    int64    `json:"userProfilePK"`
	GoalCategory int      `json:"userGoalCategoryPK"`
	GoalType     GoalType `json:"userGoalTypePK"`
	Start        Date     `json:"startDate"`
	End          Date     `json:"endDate,omitempty"`
	Value        int      `json:"goalValue"`
	Created      Date     `json:"createDate"`
}

// GoalType represents different types of goals.
type GoalType int

// String implements Stringer.
func (t GoalType) String() string {
	switch t {
	case 0:
		return "steps-per-day"
	case 4:
		return "weight"
	case 7:
		return "floors-ascended"
	default:
		return fmt.Sprintf("unknown:%d", t)
	}
}

// Goals lists all goals for displayName of type goalType. If displayName is
// empty, the currently authenticated user will be used.
func (c *Client) Goals(displayName string, goalType int) ([]Goal, error) {
	if displayName == "" && c.Profile == nil {
		return nil, ErrNotAuthenticated
	}

	if displayName == "" && c.Profile != nil {
		displayName = c.Profile.DisplayName
	}

	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/wellness-service/wellness/wellness-goals/%s?userGoalType=%d",
		displayName,
		goalType,
	)

	goals := make([]Goal, 0, 20)

	err := c.getJSON(URL, &goals)
	if err != nil {
		return nil, err
	}

	return goals, nil
}

// AddGoal will add a new goal. If displayName is empty, the currently
// authenticated user will be used.
func (c *Client) AddGoal(displayName string, goal Goal) error {
	if displayName == "" && c.Profile == nil {
		return ErrNotAuthenticated
	}

	if displayName == "" && c.Profile != nil {
		displayName = c.Profile.DisplayName
	}

	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/wellness-service/wellness/wellness-goals/%s",
		displayName,
	)

	return c.write("POST", URL, goal, 204)
}

// DeleteGoal will delete an existing goal. If displayName is empty, the
// currently authenticated user will be used.
func (c *Client) DeleteGoal(displayName string, goalID int) error {
	if displayName == "" && c.Profile == nil {
		return ErrNotAuthenticated
	}

	if displayName == "" && c.Profile != nil {
		displayName = c.Profile.DisplayName
	}

	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/wellness-service/wellness/wellness-goals/%d/%s",
		goalID,
		displayName,
	)

	return c.write("DELETE", URL, nil, 204)
}

// UpdateGoal will update an existing goal.
func (c *Client) UpdateGoal(displayName string, goal Goal) error {
	if displayName == "" && c.Profile == nil {
		return ErrNotAuthenticated
	}

	if displayName == "" && c.Profile != nil {
		displayName = c.Profile.DisplayName
	}

	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/wellness-service/wellness/wellness-goals/%d/%s",
		goal.ID,
		displayName,
	)

	return c.write("PUT", URL, goal, 204)
}
