package connect

// LastUsed describes the last synchronization.
type LastUsed struct {
	DeviceID             int    `json:"userDeviceId"`
	ProfileNumber        int    `json:"userProfileNumber"`
	ApplicationNumber    int    `json:"applicationNumber"`
	DeviceApplicationKey string `json:"lastUsedDeviceApplicationKey"`
	DeviceName           string `json:"lastUsedDeviceName"`
	DeviceUploadTime     Time   `json:"lastUsedDeviceUploadTime"`
	ImageURL             string `json:"imageUrl"`
	Released             bool   `json:"released"`
}

// LastUsed will return information about the latest synchronization.
func (c *Client) LastUsed(displayName string) (*LastUsed, error) {
	URL := "https://connect.garmin.com/modern/proxy/device-service/deviceservice/userlastused/" + displayName

	lastused := new(LastUsed)

	err := c.getJSON(URL, lastused)
	if err != nil {
		return nil, err
	}

	return lastused, err
}
