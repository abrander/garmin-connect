package connect

// Everything from https://connect.garmin.com/modern/proxy/badge-service/badge/attributes

type BadgeType struct {
	ID  int    `json:"badgeTypeId"`
	Key string `json:"badgeTypeKey"`
}

type BadgeCategory struct {
	ID  int    `json:"badgeCategoryId"`
	Key string `json:"badgeCategoryKey"`
}

type BadgeDifficulty struct {
	ID     int    `json:"badgeDifficultyId"`
	Key    string `json:"badgeDifficultyKey"`
	Points int    `json:"badgePoints"`
}

type BadgeUnit struct {
	ID  int    `json:"badgeUnitId"`
	Key string `json:"badgeUnitKey"`
}

type BadgeAssocType struct {
	ID  int    `json:"badgeAssocTypeId"`
	Key string `json:"badgeAssocTypeKey"`
}

type BadgeAttributes struct {
	BadgeTypes        []BadgeType       `json:"badgeTypes"`
	BadgeCategories   []BadgeCategory   `json:"badgeCategories"`
	BadgeDifficulties []BadgeDifficulty `json:"badgeDifficulties"`
	BadgeUnits        []BadgeUnit       `json:"badgeUnits"`
	BadgeAssocTypes   []BadgeAssocType  `json:"badgeAssocTypes"`
}

// BadgeAttributes retrieves a list of badge attributes. At time of writing
// we're not sure how these can be utilized.
func (c *Client) BadgeAttributes() (*BadgeAttributes, error) {
	URL := "https://connect.garmin.com/modern/proxy/badge-service/badge/attributes"

	attributes := new(BadgeAttributes)

	err := c.getJSON(URL, &attributes)
	if err != nil {
		return nil, err
	}

	return attributes, nil
}
