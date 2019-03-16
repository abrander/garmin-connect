package connect

import (
	"fmt"
)

// GroupAnnouncement describes a group announcement. Only one announcement can
// exist per group.
type GroupAnnouncement struct {
	ID               int    `json:"announcementId"`
	GroupID          int    `json:"groupId"`
	Title            string `json:"title"`
	Message          string `json:"message"`
	ExpireDate       Time   `json:"expireDate"`
	AnnouncementDate Time   `json:"announcementDate"`
}

// GroupAnnouncement returns the announcement for groupID.
func (c *Client) GroupAnnouncement(groupID int) (*GroupAnnouncement, error) {
	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/group-service/group/%d/announcement",
		groupID,
	)

	announcement := new(GroupAnnouncement)
	err := c.getJSON(URL, announcement)
	if err != nil {
		return nil, err
	}

	return announcement, nil
}
