package connect

// BiometricProfile holds key biometric data.
type BiometricProfile struct {
	UserID        int     `json:"userId"`
	Height        float64 `json:"height"`
	Weight        float64 `json:"weight"` // grams
	VO2Max        float64 `json:"vo2Max"`
	VO2MaxCycling float64 `json:"vo2MaxCycling"`
}

// UserInfo is very basic information about a user.
type UserInfo struct {
	Gender   string `json:"genderType"`
	Email    string `json:"email"`
	Locale   string `json:"locale"`
	TimeZone string `json:"timezone"`
	Age      int    `json:"age"`
}

// PersonalInformation is user info and a biometric profile for a user.
type PersonalInformation struct {
	UserInfo         UserInfo         `json:"userInfo"`
	BiometricProfile BiometricProfile `json:"biometricProfile"`
}

// PersonalInformation will retrieve personal information for displayName.
func (c *Client) PersonalInformation(displayName string) (*PersonalInformation, error) {
	URL := "https://connect.garmin.com/modern/proxy/userprofile-service/userprofile/personal-information/" + displayName

	pi := new(PersonalInformation)

	err := c.getJSON(URL, pi)
	if err != nil {
		return nil, err
	}

	return pi, nil
}
