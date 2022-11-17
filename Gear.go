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
	ImageNameLarge  string  `json:"imageNameLarge"`
	ImageNameMedium string  `json:"imageNameMedium"`
	ImageNameSmall  string  `json:"imageNameSmall"`
	DateBegin       Time    `json:"dateBegin"`
	DateEnd         Time    `json:"dateEnd"`
	MaximumMeters   float64 `json:"maximumMeters"`
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

// Gear will retrieve the details of the users gear
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

// GearType will list the gear types
func (c *Client) GearType() ([]GearType, error) {
	URL := "https://connect.garmin.com/modern/proxy/gear-service/gear/types"
	var gearType []GearType
	err := c.getJSON(URL, &gearType)
	if err != nil {
		return nil, err
	}

	return gearType, nil
}

// GearStats will get the statistics of an item of gear, given the uuid
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

// GearLink will link an item of gear to an activity. Multiple items of gear can be linked.
func (c *Client) GearLink(uuid string, activityID int) error {
	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/gear-service/gear/link/%s/activity/%d",
		uuid,
		activityID,
	)

	return c.write("PUT", URL, "", 200)
}

// GearUnlink will remove an item of gear from an activity. All items of gear can be unlinked.
func (c *Client) GearUnlink(uuid string, activityID int) error {
	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/gear-service/gear/unlink/%s/activity/%d",
		uuid,
		activityID,
	)

	return c.write("PUT", URL, "", 200)
}

// GearForActivity will retrieve the gear associated with an activity
func (c *Client) GearForActivity(profileID int64, activityID int) ([]Gear, error) {
	if profileID == 0 && c.Profile == nil {
		return nil, ErrNotAuthenticated
	}

	if profileID == 0 && c.Profile != nil {
		profileID = c.Profile.ProfileID
	}

	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/gear-service/gear/filterGear?userProfilePk=%d&activityId=%d",
		profileID, activityID,
	)
	var gear []Gear
	err := c.getJSON(URL, &gear)
	if err != nil {
		return nil, err
	}

	return gear, nil
}
