package connect

import (
	"fmt"
)

// Badge describes a badge.
type Badge struct {
	ID                 int     `json:"badgeId"`
	Key                string  `json:"badgeKey"`
	Name               string  `json:"badgeName"`
	CategoryID         int     `json:"badgeCategoryId"`
	DifficultyID       int     `json:"badgeDifficultyId"`
	Points             int     `json:"badgePoints"`
	TypeID             []int   `json:"badgeTypeIds"`
	SeriesID           int     `json:"badgeSeriesId"`
	Start              Time    `json:"badgeStartDate"`
	End                Time    `json:"badgeEndDate"`
	UserProfileID      int     `json:"userProfileId"`
	FullName           string  `json:"fullName"`
	DisplayName        string  `json:"displayName"`
	EarnedDate         Time    `json:"badgeEarnedDate"`
	EarnedNumber       int     `json:"badgeEarnedNumber"`
	Viewed             bool    `json:"badgeIsViewed"`
	Progress           float64 `json:"badgeProgressValue"`
	Target             float64 `json:"badgeTargetValue"`
	UnitID             int     `json:"badgeUnitId"`
	BadgeAssocTypeID   int     `json:"badgeAssocTypeId"`
	BadgeAssocDataID   string  `json:"badgeAssocDataId"`
	BadgeAssocDataName string  `json:"badgeAssocDataName"`
	EarnedByMe         bool    `json:"earnedByMe"`
	RelatedBadges      []Badge `json:"relatedBadges"`
	Connections        []Badge `json:"connections"`
}

// BadgeDetail will return details about a badge.
func (c *Client) BadgeDetail(badgeID int) (*Badge, error) {
	// Alternative URL:
	// https://connect.garmin.com/modern/proxy/badge-service/badge/DISPLAYNAME/earned/detail/BADGEID
	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/badge-service/badge/detail/v2/%d",
		badgeID)

	badge := new(Badge)

	err := c.getJSON(URL, badge)

	// This is interesting. Garmin returns 400 if an unknown badge is
	// requested. We have no way of detecting that, so we silently changes
	// the error to ErrNotFound.
	if err == ErrBadRequest {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return badge, nil
}
