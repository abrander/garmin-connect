package connect

import (
	"fmt"
)

// ActivityWeather describes the weather during an activity.
type ActivityWeather struct {
	Temperature               int     `json:"temp"`
	ApparentTemperature       int     `json:"apparentTemp"`
	DewPoint                  int     `json:"dewPoint"`
	RelativeHumidity          int     `json:"relativeHumidity"`
	WindDirection             int     `json:"windDirection"`
	WindDirectionCompassPoint string  `json:"windDirectionCompassPoint"`
	WindSpeed                 int     `json:"windSpeed"`
	Latitude                  float64 `json:"latitude"`
	Longitude                 float64 `json:"longitude"`
}

// ActivityWeather returns the reported weather for an activity.
func (c *Client) ActivityWeather(activityID int) (*ActivityWeather, error) {
	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/weather-service/weather/%d",
		activityID,
	)

	weather := new(ActivityWeather)

	err := c.getJSON(URL, weather)
	if err != nil {
		return nil, err
	}

	return weather, nil
}
