package connect

import (
	"fmt"
	"time"
)

// DateValue is a numeric value recorded on a given date.
type DateValue struct {
	Date  Date    `json:"calendarDate"`
	Value float64 `json:"value"`
}

// DailySummaries provides a daily summary of various statistics for multiple
// days.
type DailySummaries struct {
	Start                    time.Time   `json:"statisticsStartDate"`
	End                      time.Time   `json:"statisticsEndDate"`
	TotalSteps               []DateValue `json:"WELLNESS_TOTAL_STEPS"`
	ActiveCalories           []DateValue `json:"COMMON_ACTIVE_CALORIES"`
	FloorsAscended           []DateValue `json:"WELLNESS_FLOORS_ASCENDED"`
	IntensityMinutes         []DateValue `json:"WELLNESS_USER_INTENSITY_MINUTES_GOAL"`
	MaxHeartRate             []DateValue `json:"WELLNESS_MAX_HEART_RATE"`
	MinimumAverageHeartRate  []DateValue `json:"WELLNESS_MIN_AVG_HEART_RATE"`
	MinimumHeartrate         []DateValue `json:"WELLNESS_MIN_HEART_RATE"`
	AverageStress            []DateValue `json:"WELLNESS_AVERAGE_STRESS"`
	RestingHeartRate         []DateValue `json:"WELLNESS_RESTING_HEART_RATE"`
	MaxStress                []DateValue `json:"WELLNESS_MAX_STRESS"`
	AbnormalHeartRateAlers   []DateValue `json:"WELLNESS_ABNORMALHR_ALERTS_COUNT"`
	MaximumAverageHeartRate  []DateValue `json:"WELLNESS_MAX_AVG_HEART_RATE"`
	StepGoal                 []DateValue `json:"WELLNESS_TOTAL_STEP_GOAL"`
	FlorsAscendedGoal        []DateValue `json:"WELLNESS_USER_FLOORS_ASCENDED_GOAL"`
	ModerateIntensityMinutes []DateValue `json:"WELLNESS_MODERATE_INTENSITY_MINUTES"`
	TotalColaries            []DateValue `json:"WELLNESS_TOTAL_CALORIES"`
	BodyBatteryCharged       []DateValue `json:"WELLNESS_BODYBATTERY_CHARGED"`
	FloorsDescended          []DateValue `json:"WELLNESS_FLOORS_DESCENDED"`
	BMRCalories              []DateValue `json:"WELLNESS_BMR_CALORIES"`
	FoodCaloriesRemainin     []DateValue `json:"FOOD_CALORIES_REMAINING"`
	TotalCalories            []DateValue `json:"COMMON_TOTAL_CALORIES"`
	BodyBatteryDrained       []DateValue `json:"WELLNESS_BODYBATTERY_DRAINED"`
	AverageSteps             []DateValue `json:"WELLNESS_AVERAGE_STEPS"`
	VigorousIntensifyMinutes []DateValue `json:"WELLNESS_VIGOROUS_INTENSITY_MINUTES"`
	WellnessDistance         []DateValue `json:"WELLNESS_TOTAL_DISTANCE"`
	Distance                 []DateValue `json:"COMMON_TOTAL_DISTANCE"`
	WellnessActiveCalories   []DateValue `json:"WELLNESS_ACTIVE_CALORIES"`
}

// DailySummary is an extensive summary for a single day.
type DailySummary struct {
	ProfileID                        int64         `json:"userProfileId"`
	TotalKilocalories                float64       `json:"totalKilocalories"`
	ActiveKilocalories               float64       `json:"activeKilocalories"`
	BMRKilocalories                  float64       `json:"bmrKilocalories"`
	WellnessKilocalories             float64       `json:"wellnessKilocalories"`
	BurnedKilocalories               float64       `json:"burnedKilocalories"`
	ConsumedKilocalories             float64       `json:"consumedKilocalories"`
	RemainingKilocalories            float64       `json:"remainingKilocalories"`
	TotalSteps                       int           `json:"totalSteps"`
	NetCalorieGoal                   float64       `json:"netCalorieGoal"`
	TotalDistanceMeters              int           `json:"totalDistanceMeters"`
	WellnessDistanceMeters           int           `json:"wellnessDistanceMeters"`
	WellnessActiveKilocalories       float64       `json:"wellnessActiveKilocalories"`
	NetRemainingKilocalories         float64       `json:"netRemainingKilocalories"`
	UserID                           int64         `json:"userDailySummaryId"`
	Date                             Date          `json:"calendarDate"`
	UUID                             string        `json:"uuid"`
	StepGoal                         int           `json:"dailyStepGoal"`
	StartTimeGMT                     Time          `json:"wellnessStartTimeGmt"`
	EndTimeGMT                       Time          `json:"wellnessEndTimeGmt"`
	StartLocal                       Time          `json:"wellnessStartTimeLocal"`
	EndLocal                         Time          `json:"wellnessEndTimeLocal"`
	Duration                         time.Duration `json:"durationInMilliseconds"`
	Description                      string        `json:"wellnessDescription"`
	HighlyActive                     time.Duration `json:"highlyActiveSeconds"`
	Active                           time.Duration `json:"activeSeconds"`
	Sedentary                        time.Duration `json:"sedentarySeconds"`
	Sleeping                         time.Duration `json:"sleepingSeconds"`
	IncludesWellnessData             bool          `json:"includesWellnessData"`
	IncludesActivityData             bool          `json:"includesActivityData"`
	IncludesCalorieConsumedData      bool          `json:"includesCalorieConsumedData"`
	PrivacyProtected                 bool          `json:"privacyProtected"`
	ModerateIntensity                time.Duration `json:"moderateIntensityMinutes"`
	VigorousIntensity                time.Duration `json:"vigorousIntensityMinutes"`
	FloorsAscendedInMeters           float64       `json:"floorsAscendedInMeters"`
	FloorsDescendedInMeters          float64       `json:"floorsDescendedInMeters"`
	FloorsAscended                   float64       `json:"floorsAscended"`
	FloorsDescended                  float64       `json:"floorsDescended"`
	IntensityGoal                    time.Duration `json:"intensityMinutesGoal"`
	FloorsAscendedGoal               int           `json:"userFloorsAscendedGoal"`
	MinHeartRate                     int           `json:"minHeartRate"`
	MaxHeartRate                     int           `json:"maxHeartRate"`
	RestingHeartRate                 int           `json:"restingHeartRate"`
	LastSevenDaysAvgRestingHeartRate int           `json:"lastSevenDaysAvgRestingHeartRate"`
	Source                           string        `json:"source"`
	AverageStress                    int           `json:"averageStressLevel"`
	MaxStress                        int           `json:"maxStressLevel"`
	Stress                           time.Duration `json:"stressDuration"`
	RestStress                       time.Duration `json:"restStressDuration"`
	ActivityStress                   time.Duration `json:"activityStressDuration"`
	UncategorizedStress              time.Duration `json:"uncategorizedStressDuration"`
	TotalStress                      time.Duration `json:"totalStressDuration"`
	LowStress                        time.Duration `json:"lowStressDuration"`
	MediumStress                     time.Duration `json:"mediumStressDuration"`
	HighStress                       time.Duration `json:"highStressDuration"`
	StressQualifier                  string        `json:"stressQualifier"`
	MeasurableAwake                  time.Duration `json:"measurableAwakeDuration"`
	MeasurableAsleep                 time.Duration `json:"measurableAsleepDuration"`
	LastSyncGMT                      Time          `json:"lastSyncTimestampGMT"`
	MinAverageHeartRate              int           `json:"minAvgHeartRate"`
	MaxAverageHeartRate              int           `json:"maxAvgHeartRate"`
}

// DailySummary will retrieve a detailed daily summary for date. If
// displayName is empty, the currently authenticated user will be used.
func (c *Client) DailySummary(displayName string, date time.Time) (*DailySummary, error) {
	if displayName == "" && c.Profile == nil {
		return nil, ErrNotAuthenticated
	}

	if displayName == "" && c.Profile != nil {
		displayName = c.Profile.DisplayName
	}

	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/usersummary-service/usersummary/daily/%s?calendarDate=%s",
		displayName,
		formatDate(date),
	)

	summary := new(DailySummary)

	err := c.getJSON(URL, summary)
	if err != nil {
		return nil, err
	}

	summary.Duration *= time.Millisecond
	summary.HighlyActive *= time.Second
	summary.Active *= time.Second
	summary.Sedentary *= time.Second
	summary.Sleeping *= time.Second
	summary.ModerateIntensity *= time.Minute
	summary.VigorousIntensity *= time.Minute
	summary.IntensityGoal *= time.Minute
	summary.Stress *= time.Second
	summary.RestStress *= time.Second
	summary.ActivityStress *= time.Second
	summary.UncategorizedStress *= time.Second
	summary.TotalStress *= time.Second
	summary.LowStress *= time.Second
	summary.MediumStress *= time.Second
	summary.HighStress *= time.Second
	summary.MeasurableAwake *= time.Second
	summary.MeasurableAsleep *= time.Second

	return summary, nil
}

// DailySummaries will retrieve a daily summary for userID.
func (c *Client) DailySummaries(userID string, from time.Time, until time.Time) (*DailySummaries, error) {
	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/userstats-service/wellness/daily/%s?fromDate=%s&untilDate=%s",
		userID,
		formatDate(from),
		formatDate(until),
	)

	if !c.authenticated() {
		return nil, ErrNotAuthenticated
	}

	// We use a proxy object to deserialize the values to proper Go types.
	var proxy struct {
		Start      Date `json:"statisticsStartDate"`
		End        Date `json:"statisticsEndDate"`
		AllMetrics struct {
			Summary DailySummaries `json:"metricsMap"`
		} `json:"allMetrics"`
	}

	err := c.getJSON(URL, &proxy)
	if err != nil {
		return nil, err
	}

	ret := &proxy.AllMetrics.Summary
	ret.Start = proxy.Start.Time()
	ret.End = proxy.End.Time()

	return ret, nil
}
