package connect

// SocialProfile represents a Garmin Connect user.
type SocialProfile struct {
	ID                    int64    `json:"id"`
	ProfileID             int64    `json:"profileId"`
	ConnectionRequestID   int      `json:"connectionRequestId"`
	GarminGUID            string   `json:"garminGUID"`
	DisplayName           string   `json:"displayName"`
	Fullname              string   `json:"fullName"`
	Username              string   `json:"userName"`
	ProfileImageURLLarge  string   `json:"profileImageUrlLarge"`
	ProfileImageURLMedium string   `json:"profileImageUrlMedium"`
	ProfileImageURLSmall  string   `json:"profileImageUrlSmall"`
	Location              string   `json:"location"`
	FavoriteActivityTypes []string `json:"favoriteActivityTypes"`
	UserRoles             []string `json:"userRoles"`
	UserProfileFullName   string   `json:"userProfileFullName"`
	UserLevel             int      `json:"userLevel"`
	UserPoint             int      `json:"userPoint"`
}

// SocialProfile retrieves a profile for a Garmin Connect user. If displayName
// is empty, the profile for the currently authenticated user will be returned.
func (c *Client) SocialProfile(displayName string) (*SocialProfile, error) {
	URL := "https://connect.garmin.com/modern/proxy/userprofile-service/socialProfile/" + displayName

	profile := new(SocialProfile)

	err := c.getJSON(URL, profile)
	if err != nil {
		return nil, err
	}

	return profile, err
}

// PublicSocialProfile retrieves the public profile for displayName.
func (c *Client) PublicSocialProfile(displayName string) (*SocialProfile, error) {
	URL := "https://connect.garmin.com/modern/proxy/userprofile-service/socialProfile/public/" + displayName

	profile := new(SocialProfile)

	err := c.getJSON(URL, profile)
	if err != nil {
		return nil, err
	}

	return profile, err
}

// BlockedUsers returns the list of blocked users for the currently
// authenticated user.
func (c *Client) BlockedUsers() ([]SocialProfile, error) {
	URL := "https://connect.garmin.com/modern/proxy/userblock-service/blockuser"

	var results []SocialProfile

	err := c.getJSON(URL, &results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// BlockUser will block a user.
func (c *Client) BlockUser(displayName string) error {
	URL := "https://connect.garmin.com/modern/proxy/userblock-service/blockuser/" + displayName

	return c.write("POST", URL, nil, 200)
}

// UnblockUser removed displayName from the block list.
func (c *Client) UnblockUser(displayName string) error {
	URL := "https://connect.garmin.com/modern/proxy/userblock-service/blockuser/" + displayName

	return c.write("DELETE", URL, nil, 204)
}
