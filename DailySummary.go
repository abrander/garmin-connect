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

// DailySummary provides a daily summary of various statistics.
type DailySummary struct {
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

// DailySummary will retrieve a daily summary for userID.
func (c *Client) DailySummary(userID string, from time.Time, until time.Time) (*DailySummary, error) {
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
			Summary DailySummary `json:"metricsMap"`
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
