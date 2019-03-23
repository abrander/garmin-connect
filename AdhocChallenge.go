package connect

import (
	"fmt"
)

// Player represents a participant in a challenge.
type Player struct {
	UserProfileID         int     `json:"userProfileId"`
	TotalNumber           float64 `json:"totalNumber"`
	LastSyncTime          Time    `json:"lastSyncTime"`
	Ranking               int     `json:"ranking"`
	ProfileImageURLSmall  string  `json:"profileImageSmall"`
	ProfileImageURLMedium string  `json:"profileImageMedium"`
	FullName              string  `json:"fullName"`
	DisplayName           string  `json:"displayName"`
	ProUser               bool    `json:"isProUser"`
	TodayNumber           float64 `json:"todayNumber"`
	AcceptedChallenge     bool    `json:"isAcceptedChallenge"`
}

// AdhocChallenge is a user-initiated challenge between 2 or more participants.
type AdhocChallenge struct {
	SocialChallengeStatusID       int      `json:"socialChallengeStatusId"`
	SocialChallengeActivityTypeID int      `json:"socialChallengeActivityTypeId"`
	SocialChallengeType           int      `json:"socialChallengeType"`
	Name                          string   `json:"adHocChallengeName"`
	Description                   string   `json:"adHocChallengeDesc"`
	OwnerProfileID                int      `json:"ownerUserProfileId"`
	UUID                          string   `json:"uuid"`
	Start                         Time     `json:"startDate"`
	End                           Time     `json:"endDate"`
	DurationTypeID                int      `json:"durationTypeId"`
	UserRanking                   int      `json:"userRanking"`
	Players                       []Player `json:"players"`
}

// AdhocChallenges will list the currently non-completed Ad-Hoc challenges.
// Please note that Players will not be populated, use AdhocChallenge() to
// retrieve players for a challenge.
func (c *Client) AdhocChallenges() ([]AdhocChallenge, error) {
	URL := "https://connect.garmin.com/modern/proxy/adhocchallenge-service/adHocChallenge/nonCompleted"

	if !c.authenticated() {
		return nil, ErrNotAuthenticated
	}

	challenges := make([]AdhocChallenge, 0, 10)

	err := c.getJSON(URL, &challenges)
	if err != nil {
		return nil, err
	}

	return challenges, nil
}

// HistoricalAdhocChallenges will retrieve the list of completed ad-hoc
// challenges.
func (c *Client) HistoricalAdhocChallenges() ([]AdhocChallenge, error) {
	URL := "https://connect.garmin.com/modern/proxy/adhocchallenge-service/adHocChallenge/historical"

	if !c.authenticated() {
		return nil, ErrNotAuthenticated
	}

	challenges := make([]AdhocChallenge, 0, 100)

	err := c.getJSON(URL, &challenges)
	if err != nil {
		return nil, err
	}

	return challenges, nil
}

// AdhocChallenge will retrieve details for challenge with uuid.
func (c *Client) AdhocChallenge(uuid string) (*AdhocChallenge, error) {
	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/adhocchallenge-service/adHocChallenge/%s", uuid)

	challenge := new(AdhocChallenge)

	err := c.getJSON(URL, challenge)
	if err != nil {
		return nil, err
	}

	return challenge, nil
}

// LeaveAdhocChallenge will leave an ad-hoc challenge. If profileID is 0, the
// currently authenticated user will be used.
func (c *Client) LeaveAdhocChallenge(challengeUUID string, profileID int64) error {
	if profileID == 0 && c.Profile == nil {
		return ErrNotAuthenticated
	}

	if profileID == 0 && c.Profile != nil {
		profileID = c.Profile.ProfileID
	}

	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/adhocchallenge-service/adHocChallenge/%s/player/%d",
		challengeUUID,
		profileID,
	)

	return c.write("DELETE", URL, nil, 0)
}
