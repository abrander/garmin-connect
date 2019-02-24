package connect

// SocialProfile represents a Garmin Connect user.
type SocialProfile struct {
	ID                    int      `json:"id"`
	ProfileID             int      `json:"profileId"`
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
