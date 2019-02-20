package connect

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
}
