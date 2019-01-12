package connect

// Timezones is the list of known time zones in Garmin Connect.
type Timezones []Timezone

// Timezones will retrieve the list of known timezones in Garmin Connect.
func (c *Client) Timezones() (Timezones, error) {
	URL := "https://connect.garmin.com/modern/proxy/system-service/timezoneUnits"

	if !c.authenticated() {
		return nil, ErrNotAuthenticated
	}

	timezones := make(Timezones, 0, 100)

	err := c.getJSON(URL, &timezones)
	if err != nil {
		return nil, err
	}

	return timezones, nil
}

// FindID will search for the timezone with id.
func (ts Timezones) FindID(id int) (Timezone, bool) {
	for _, t := range ts {
		if t.ID == id {
			return t, true
		}
	}

	return Timezone{}, false
}

// FindKey will search for the timezone with key key.
func (ts Timezones) FindKey(key string) (Timezone, bool) {
	for _, t := range ts {
		if t.Key == key {
			return t, true
		}
	}

	return Timezone{}, false
}
