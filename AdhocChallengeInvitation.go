package connect

import (
	"fmt"
)

// AdhocChallengeInvitation is a ad-hoc challenge invitation.
type AdhocChallengeInvitation struct {
	AdhocChallenge `json:",inline"`

	UUID               string `json:"adHocChallengeUuid"`
	InviteID           int    `json:"adHocChallengeInviteId"`
	InvitorName        string `json:"invitorName"`
	InvitorID          int    `json:"invitorId"`
	InvitorDisplayName string `json:"invitorDisplayName"`
	InviteeID          int    `json:"inviteeId"`
	UserImageURL       string `json:"userImageUrl"`
}

// AdhocChallengeInvites list Ad-Hoc challenges awaiting response.
func (c *Client) AdhocChallengeInvites() ([]AdhocChallengeInvitation, error) {
	URL := "https://connect.garmin.com/modern/proxy/adhocchallenge-service/adHocChallenge/invite"

	if !c.authenticated() {
		return nil, ErrNotAuthenticated
	}

	challenges := make([]AdhocChallengeInvitation, 0, 10)

	err := c.getJSON(URL, &challenges)
	if err != nil {
		return nil, err
	}

	// Make sure the embedded UUID matches in case the user uses the embedded
	// AdhocChallenge for something.
	for i := range challenges {
		challenges[i].AdhocChallenge.UUID = challenges[i].UUID
	}

	return challenges, nil
}

// AdhocChallengeInvitationRespond will respond to a ad-hoc challenge. If
// accept is false, the challenge will be declined.
func (c *Client) AdhocChallengeInvitationRespond(inviteID int, accept bool) error {
	scope := "decline"
	if accept {
		scope = "accept"
	}

	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/adhocchallenge-service/adHocChallenge/invite/%d/%s", inviteID, scope)

	payload := struct {
		InviteID int    `json:"inviteId"`
		Scope    string `json:"scope"`
	}{
		inviteID,
		scope,
	}

	return c.write("PUT", URL, payload, 0)
}
