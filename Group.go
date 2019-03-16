package connect

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Group describes a Garmin Connect group.
type Group struct {
	ID                    int            `json:"id"`
	Name                  string         `json:"groupName"`
	Description           string         `json:"groupDescription"`
	OwnerID               int            `json:"ownerId"`
	ProfileImageURLLarge  string         `json:"profileImageUrlLarge"`
	ProfileImageURLMedium string         `json:"profileImageUrlMedium"`
	ProfileImageURLSmall  string         `json:"profileImageUrlSmall"`
	Visibility            string         `json:"groupVisibility"`
	Privacy               string         `json:"groupPrivacy"`
	Location              string         `json:"location"`
	WebsiteURL            string         `json:"websiteUrl"`
	FacebookURL           string         `json:"facebookUrl"`
	TwitterURL            string         `json:"twitterUrl"`
	PrimaryActivities     []string       `json:"primaryActivities"`
	OtherPrimaryActivity  string         `json:"otherPrimaryActivity"`
	LeaderboardTypes      []string       `json:"leaderboardTypes"`
	FeatureTypes          []string       `json:"featureTypes"`
	CorporateWellness     bool           `json:"isCorporateWellness"`
	ActivityFeedTypes     []ActivityType `json:"activityFeedTypes"`
}

/*
Unknowns:
"membershipStatus": null,
"isCorporateWellness": false,
"programName": null,
"programTextColor": null,
"programBackgroundColor": null,
"groupMemberCount": null,
*/

// Groups will return the group membership. If displayName is empty, the
// currently authenticated user will be used.
func (c *Client) Groups(displayName string) ([]Group, error) {
	if displayName == "" && c.Profile == nil {
		return nil, ErrNotAuthenticated
	}

	if displayName == "" && c.Profile != nil {
		displayName = c.Profile.DisplayName
	}

	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/group-service/groups/%s", displayName)

	groups := make([]Group, 0, 30)

	err := c.getJSON(URL, &groups)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

// SearchGroups can search for groups in Garmin Connect.
func (c *Client) SearchGroups(keyword string) ([]Group, error) {
	URL := "https://connect.garmin.com/modern/proxy/group-service/keyword"

	payload := url.Values{
		"start":   {"1"},
		"limit":   {"20"},
		"keyword": {keyword},
	}

	req, err := c.newRequest("POST", URL, strings.NewReader(payload.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("nk", "NT")
	req.Header.Add("content-type", "application/x-www-form-urlencoded; charset=UTF-8")

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var proxy struct {
		Groups []Group `json:"groupDTOs"`
	}

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&proxy)
	if err != nil {
		return nil, err
	}

	return proxy.Groups, nil
}

// Group returns details about groupID.
func (c *Client) Group(groupID int) (*Group, error) {
	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/group-service/group/%d", groupID)

	group := new(Group)

	err := c.getJSON(URL, group)
	if err != nil {
		return nil, err
	}

	return group, nil
}

// JoinGroup joins a group. If profileID is 0, the currently authenticated
// user will be used.
func (c *Client) JoinGroup(groupID int) error {
	if c.Profile == nil {
		return ErrNotAuthenticated
	}

	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/group-service/group/%d/member/%d",
		groupID,
		c.Profile.ProfileID,
	)

	payload := struct {
		GroupID   int     `json:"groupId"`
		Role      *string `json:"groupRole"` // is always null?
		ProfileID int     `json:"userProfileId"`
	}{
		groupID,
		nil,
		c.Profile.ProfileID,
	}

	body := bytes.NewBuffer(nil)
	enc := json.NewEncoder(body)
	err := enc.Encode(payload)
	if err != nil {
		return err
	}

	req, err := c.newRequest("POST", URL, body)
	if err != nil {
		return err
	}

	req.Header.Add("nk", "NT")
	req.Header.Add("content-type", "application/json")

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected status code %d", resp.StatusCode)
	}

	return nil
}

// LeaveGroup leaves a group.
func (c *Client) LeaveGroup(groupID int) error {
	if c.Profile == nil {
		return ErrNotAuthenticated
	}

	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/group-service/group/%d/member/%d",
		groupID,
		c.Profile.ProfileID,
	)

	req, err := c.newRequest("DELETE", URL, nil)
	if err != nil {
		return err
	}

	req.Header.Add("nk", "NT")

	resp, err := c.do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		return fmt.Errorf("HTTP call returned %d", resp.StatusCode)
	}

	return nil
}
