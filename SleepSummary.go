package connect

import (
	"fmt"
	"time"
)

// "sleepQualityTypePK": null,
// "sleepResultTypePK": null,

// SleepSummary is a summary of sleep for a single night.
type SleepSummary struct {
	ID               int64         `json:"id"`
	UserProfilePK    int64         `json:"userProfilePK"`
	Sleep            time.Duration `json:"sleepTimeSeconds"`
	Nap              time.Duration `json:"napTimeSeconds"`
	Confirmed        bool          `json:"sleepWindowConfirmed"`
	Confirmation     string        `json:"sleepWindowConfirmationType"`
	StartGMT         Time          `json:"sleepStartTimestampGMT"`
	EndGMT           Time          `json:"sleepEndTimestampGMT"`
	StartLocal       Time          `json:"sleepStartTimestampLocal"`
	EndLocal         Time          `json:"sleepEndTimestampLocal"`
	AutoStartGMT     Time          `json:"autoSleepStartTimestampGMT"`
	AutoEndGMT       Time          `json:"autoSleepEndTimestampGMT"`
	Unmeasurable     time.Duration `json:"unmeasurableSleepSeconds"`
	Deep             time.Duration `json:"deepSleepSeconds"`
	Light            time.Duration `json:"lightSleepSeconds"`
	REM              time.Duration `json:"remSleepSeconds"`
	Awake            time.Duration `json:"awakeSleepSeconds"`
	DeviceRemCapable bool          `json:"deviceRemCapable"`
	REMData          bool          `json:"remData"`
}

// SleepMovement denotes the amount of movement for a short time period
// during sleep.
type SleepMovement struct {
	Start Time    `json:"startGMT"`
	End   Time    `json:"endGMT"`
	Level float64 `json:"activityLevel"`
}

// SleepLevel represents the sleep level for a longer period of time.
type SleepLevel struct {
	Start Time       `json:"startGMT"`
	End   Time       `json:"endGMT"`
	State SleepState `json:"activityLevel"`
}

// SleepData will retrieve sleep data for date for a given displayName. If
// displayName is empty, the currently authenticated user will be used.
func (c *Client) SleepData(displayName string, date time.Time) (*SleepSummary, []SleepMovement, []SleepLevel, error) {
	if displayName == "" && c.Profile == nil {
		return nil, nil, nil, ErrNotAuthenticated
	}

	if displayName == "" && c.Profile != nil {
		displayName = c.Profile.DisplayName
	}

	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/wellness-service/wellness/dailySleepData/%s?date=%s&nonSleepBufferMinutes=60",
		displayName,
		formatDate(date),
	)

	var proxy struct {
		SleepSummary SleepSummary    `json:"dailySleepDTO"`
		REMData      bool            `json:"remSleepData"`
		Movement     []SleepMovement `json:"sleepMovement"`
		Levels       []SleepLevel    `json:"sleepLevels"`
	}

	err := c.getJSON(URL, &proxy)
	if err != nil {
		return nil, nil, nil, err
	}

	// All timings from Garmin are in seconds.
	proxy.SleepSummary.Sleep *= time.Second
	proxy.SleepSummary.Nap *= time.Second
	proxy.SleepSummary.Unmeasurable *= time.Second
	proxy.SleepSummary.Deep *= time.Second
	proxy.SleepSummary.Light *= time.Second
	proxy.SleepSummary.REM *= time.Second
	proxy.SleepSummary.Awake *= time.Second

	proxy.SleepSummary.REMData = proxy.REMData

	return &proxy.SleepSummary, proxy.Movement, proxy.Levels, nil
}
