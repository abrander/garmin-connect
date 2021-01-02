package connect

import (
	"fmt"
)

// Gear describes a Garmin Connect gear entry
type Gear struct {
	Uuid            string  `json:"uuid"`
	GearPk          int     `json:"gearPk"`
	UserProfileID   int64   `json:"userProfilePk"`
	GearMakeName    string  `json:"gearMakeName"`
	GearModelName   string  `json:"gearModelName"`
	GearTypeName    string  `json:"gearTypeName"`
	DisplayName     string  `json:"displayName"`
	CustomMakeModel string  `json:"customMakeModel"`
	ImageNameLarge  *string `json:"imageNameLarge"`
	ImageNameMedium *string `json:"imageNameMedium"`
	ImageNameSmall  *string `json:"imageNameSmall"`
	DateBegin       Time    `json:"dateBegin"`
	DateEnd         *Time   `json:"dateEnd"`
	MaximumMeters   float32 `json:"maximumMeters"`
	Notified        bool    `json:"notified"`
	CreateDate      Time    `json:"createDate"`
	UpdateDate      Time    `json:"updateDate"`
}

// GearType desribes the types of gear
type GearType struct {
	TypeID     int    `json:"gearTypePk"`
	TypeName   string `json:"gearTypeName"`
	CreateDate Time   `json:"createDate"`
	UpdateDate Time   `json:"updateData"`
}

// GearStats describes the stats of gear
type GearStats struct {
	TotalDistance   float64 `json:"totalDistance"`
	TotalActivities int     `json:"totalActivities"`
	Processsing     bool    `json:"processing"`
}

func (c *Client) Gear(profileID int64) ([]Gear, error) {
	if profileID == 0 && c.Profile == nil {
		return nil, ErrNotAuthenticated
	}

	if profileID == 0 && c.Profile != nil {
		profileID = c.Profile.ProfileID
	}

	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/gear-service/gear/filterGear?userProfilePk=%d",
		profileID,
	)
	var gear []Gear
	err := c.getJSON(URL, &gear)
	if err != nil {
		return nil, err
	}
	return gear, nil
}

func (c *Client) GearType() ([]GearType, error) {
	URL := "https://connect.garmin.com/modern/proxy/gear-service/gear/types"
	var gearType []GearType
	err := c.getJSON(URL, &gearType)
	if err != nil {
		return nil, err
	}
	return gearType, nil
}

func (c *Client) GearStats(uuid string) (*GearStats, error) {
	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/userstats-service/gears/%s",
		uuid,
	)
	gearStats := new(GearStats)
	err := c.getJSON(URL, &gearStats)
	if err != nil {
		return nil, err
	}
	return gearStats, nil
}
