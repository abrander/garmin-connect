package connect

// BadgeStatus is the badge status for a Connect user.
type BadgeStatus struct {
	ProfileID             int     `json:"userProfileId"`
	Fullname              string  `json:"fullName"`
	DisplayName           string  `json:"displayName"`
	ProUser               bool    `json:"userPro"`
	ProfileImageURLLarge  string  `json:"profileImageUrlLarge"`
	ProfileImageURLMedium string  `json:"profileImageUrlMedium"`
	ProfileImageURLSmall  string  `json:"profileImageUrlSmall"`
	Level                 int     `json:"userLevel"`
	LevelUpdateTime       Time    `json:"levelUpdateDate"`
	Point                 int     `json:"userPoint"`
	Badges                []Badge `json:"badges"`
}

// BadgeLeaderBoard returns the leaderboard for points for the currently
// authenticated user.
func (c *Client) BadgeLeaderBoard() ([]BadgeStatus, error) {
	URL := "https://connect.garmin.com/modern/proxy/badge-service/badge/leaderboard"

	if !c.authenticated() {
		return nil, ErrNotAuthenticated
	}

	var proxy struct {
		LeaderBoad []BadgeStatus `json:"connections"`
	}

	err := c.getJSON(URL, &proxy)
	if err != nil {
		return nil, err
	}

	return proxy.LeaderBoad, nil
}

// BadgeCompare will compare the earned badges of the currently authenticated user against displayName.
func (c *Client) BadgeCompare(displayName string) (*BadgeStatus, *BadgeStatus, error) {
	URL := "https://connect.garmin.com/modern/proxy/badge-service/badge/compare/" + displayName

	if !c.authenticated() {
		return nil, nil, ErrNotAuthenticated
	}

	var proxy struct {
		User       *BadgeStatus `json:"user"`
		Connection *BadgeStatus `json:"connection"`
	}

	err := c.getJSON(URL, &proxy)
	if err != nil {
		return nil, nil, err
	}

	return proxy.User, proxy.Connection, nil
}

// BadgesEarned will return the list of badges earned by the curently
// authenticated user.
func (c *Client) BadgesEarned() ([]Badge, error) {
	URL := "https://connect.garmin.com/modern/proxy/badge-service/badge/earned"

	if !c.authenticated() {
		return nil, ErrNotAuthenticated
	}

	badges := make([]Badge, 0, 200)
	err := c.getJSON(URL, &badges)
	if err != nil {
		return nil, err
	}

	return badges, nil
}

// BadgesAvailable will return the list of badges not yet earned by the curently
// authenticated user.
func (c *Client) BadgesAvailable() ([]Badge, error) {
	URL := "https://connect.garmin.com/modern/proxy/badge-service/badge/available"

	if !c.authenticated() {
		return nil, ErrNotAuthenticated
	}

	badges := make([]Badge, 0, 200)
	err := c.getJSON(URL, &badges)
	if err != nil {
		return nil, err
	}

	return badges, nil
}
