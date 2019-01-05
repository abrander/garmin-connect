package connect

// SocialProfile represents a Garmin Connect user.
type SocialProfile struct {
	ID          int    `json:"id"`
	ProfileID   int    `json:"profileId"`
	GarminGUID  string `json:"garminGUID"`
	DisplayName string `json:"displayName"`
	Fullname    string `json:"fullName"`
	Username    string `json:"userName"`

	ProfileImageURLLarge  string `json:"profileImageUrlLarge"`
	ProfileImageURLMedium string `json:"profileImageUrlMedium"`
	ProfileImageURLSmall  string `json:"profileImageUrlSmall"`
	Location              string `json:"location"`
	//"facebookUrl": null,
	//"twitterUrl": null,
	//"personalWebsite": null,
	//"motivation": null,
	//"bio": null,
	//"primaryActivity": null,
	FavoriteActivityTypes []string `json:"favoriteActivityTypes"`
	//"runningTrainingSpeed": 0,
	//"cyclingTrainingSpeed": 0,
	//"favoriteCyclingActivityTypes": [],
	//"cyclingClassification": null,
	//"cyclingMaxAvgPower": 0,
	//"swimmingTrainingSpeed": 0,
	//"profileVisibility": "following",
	//"activityStartVisibility": "public",
	//"activityMapVisibility": "public",
	//"courseVisibility": "public",
	//"activityHeartRateVisibility": "public",
	//"activityPowerVisibility": "public",
	//"badgeVisibility": "following",
	//"showAge": true,
	//"showWeight": true,
	//"showHeight": true,
	//"showWeightClass": false,
	//"showAgeRange": false,
	//"showGender": true,
	//"showActivityClass": false,
	//"showVO2Max": true,
	//"showPersonalRecords": true,
	//"showLast12Months": true,
	//"showLifetimeTotals": true,
	//"showUpcomingEvents": true,
	//"showRecentFavorites": true,
	//"showRecentDevice": true,
	//"showRecentGear": false,
	//"showBadges": true,
	//"otherActivity": null,
	//"otherPrimaryActivity": null,
	//"otherMotivation": null,
	UserRoles []string `json:"userRoles"`
	//"nameApproved": true,
	UserProfileFullName string `json:"userProfileFullName"`
	//"makeGolfScorecardsPrivate": true,
	//"allowGolfLiveScoring": false,
	//"allowGolfScoringByConnections": true,
	UserLevel int `json:"userLevel"`
	UserPoint int `json:"userPoint"`
	//"levelUpdateDate": "2018-11-12T06:12:11.0",
	//"levelIsViewed": false,
	//"levelPointThreshold": 140,
	//"userPro": false
}

// SocialProfile retrieves a profile for a Garmin Connect user. If displayName
// is empty, the profile for the currently authenticated user will be returned.
func (c *Client) SocialProfile(displayName string) (*SocialProfile, error) {
	URL := "https://connect.garmin.com/modern/proxy/userprofile-service/socialProfile/" + displayName

	if !c.authenticated() {
		return nil, ErrNotAuthenticated
	}

	var profile SocialProfile
	err := c.getJSON(URL, &profile)
	return &profile, err
}
