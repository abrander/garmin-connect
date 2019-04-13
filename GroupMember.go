package connect

import (
	"fmt"
	"time"
)

// GroupMember describes a member of a group.
type GroupMember struct {
	SocialProfile

	Joined time.Time `json:"joinDate"`
	Role   string    `json:"groupRole"`
}

// GroupMembers will return the member list of a group.
func (c *Client) GroupMembers(groupID int) ([]GroupMember, error) {
	type proxy struct {
		ID                    string `json:"id"`
		GroupID               int    `json:"groupId"`
		UserProfileID         int64  `json:"userProfileId"`
		DisplayName           string `json:"displayName"`
		Location              string `json:"location"`
		Joined                Date   `json:"joinDate"`
		Role                  string `json:"groupRole"`
		Name                  string `json:"fullName"`
		ProfileImageURLLarge  string `json:"profileImageLarge"`
		ProfileImageURLMedium string `json:"profileImageMedium"`
		ProfileImageURLSmall  string `json:"profileImageSmall"`
		Pro                   bool   `json:"userPro"`
		Level                 int    `json:"userLevel"`
	}
	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/group-service/group/%d/members",
		groupID,
	)

	membersProxy := make([]proxy, 0, 100)
	err := c.getJSON(URL, &membersProxy)
	if err != nil {
		return nil, err
	}

	members := make([]GroupMember, len(membersProxy))
	for i, p := range membersProxy {
		members[i].DisplayName = p.DisplayName
		members[i].ProfileID = p.UserProfileID
		members[i].DisplayName = p.DisplayName
		members[i].Location = p.Location
		members[i].Fullname = p.Name
		members[i].ProfileImageURLLarge = p.ProfileImageURLLarge
		members[i].ProfileImageURLMedium = p.ProfileImageURLMedium
		members[i].ProfileImageURLSmall = p.ProfileImageURLSmall
		members[i].UserLevel = p.Level

		members[i].Joined = p.Joined.Time()
		members[i].Role = p.Role
	}

	return members, nil
}
