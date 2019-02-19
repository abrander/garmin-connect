package connect

// BadgeStatus is the badge status for a Connect user.
type BadgeStatus struct {
	ProfileID             int    `json:"userProfileId"`
	Fullname              string `json:"fullName"`
	DisplayName           string `json:"displayName"`
	ProUser               bool   `json:"userPro"`
	ProfileImageURLLarge  string `json:"profileImageUrlLarge"`
	ProfileImageURLMedium string `json:"profileImageUrlMedium"`
	ProfileImageURLSmall  string `json:"profileImageUrlSmall"`
	Level                 int    `json:"userLevel"`
	Point                 int    `json:"userPoint"`
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
