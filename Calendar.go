package connect

import (
	"fmt"
)

// CalendarYear describes a Garmin Connect calendar year
type CalendarYear struct {
	StartDayOfJanuary int           `json:"startDayofJanuary"`
	LeapYear          bool          `json:"leapYear"`
	YearItems         []YearItem    `json:"yearItems"`
	YearSummaries     []YearSummary `json:"yearSummaries"`
}

// YearItem describes an item on a Garmin Connect calendar year
type YearItem struct {
	Date    Date `json:"date"`
	Display int  `json:"display"`
}

// YearSummary describes a per-activity-type yearly summary on a Garmin Connect calendar year
type YearSummary struct {
	ActivityTypeID     int `json:"activityTypeId"`
	NumberOfActivities int `json:"numberOfActivities"`
	TotalDistance      int `json:"totalDistance"`
	TotalDuration      int `json:"totalDuration"`
	TotalCalories      int `json:"totalCalories"`
}

// CalendarMonth describes a Garmin Conenct calendar month
type CalendarMonth struct {
	StartDayOfMonth      int            `json:"startDayOfMonth"`
	NumOfDaysInMonth     int            `json:"numOfDaysInMonth"`
	NumOfDaysInPrevMonth int            `json:"numOfDaysInPrevMonth"`
	Month                int            `json:"month"`
	Year                 int            `json:"year"`
	CalendarItems        []CalendarItem `json:"calendarItems"`
}

// CalendarWeek describes a Garmin Connect calendar week
type CalendarWeek struct {
	StartDate        Date           `json:"startDate"`
	EndDate          Date           `json:"endDate"`
	NumOfDaysInMonth int            `json:"numOfDaysInMonth"`
	CalendarItems    []CalendarItem `json:"calendarItems"`
}

// CalendarItem describes an activity displayed on a Garmin Connect calendar
type CalendarItem struct {
	ID                       int     `json:"id"`
	ItemType                 string  `json:"itemType"`
	ActivityTypeID           int     `json:"activityTypeId"`
	Title                    string  `json:"title"`
	Date                     Date    `json:"date"`
	Duration                 int     `json:"duration"`
	Distance                 int     `json:"distance"`
	Calories                 int     `json:"calories"`
	StartTimestampLocal      Time    `json:"startTimestampLocal"`
	ElapsedDuration          float64 `json:"elapsedDuration"`
	Strokes                  float64 `json:"strokes"`
	MaxSpeed                 float64 `json:"maxSpeed"`
	ShareableEvent           bool    `json:"shareableEvent"`
	AutoCalcCalories         bool    `json:"autoCalcCalories"`
	ProtectedWorkoutSchedule bool    `json:"protectedWorkoutSchedule"`
	IsParent                 bool    `json:"isParent"`
}

// CalendarYear will get the activity summaries  and list of days active for a given year
func (c *Client) CalendarYear(year int) (*CalendarYear, error) {
	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/calendar-service/year/%d",
		year,
	)
	calendarYear := new(CalendarYear)
	err := c.getJSON(URL, &calendarYear)
	if err != nil {
		return nil, err
	}

	return calendarYear, nil
}

// CalendarMonth will get the activities for a given month
func (c *Client) CalendarMonth(year int, month int) (*CalendarMonth, error) {
	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/calendar-service/year/%d/month/%d",
		year,
		month-1, // Months in Garmin Connect start from zero
	)
	calendarMonth := new(CalendarMonth)
	err := c.getJSON(URL, &calendarMonth)
	if err != nil {
		return nil, err
	}

	return calendarMonth, nil
}

// CalendarWeek will get the activities for a given week. A week will be returned that contains the day requested, not starting with)
func (c *Client) CalendarWeek(year int, month int, week int) (*CalendarWeek, error) {
	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/calendar-service/year/%d/month/%d/day/%d/start/1",
		year,
		month-1, // Months in Garmin Connect start from zero
		week,
	)
	calendarWeek := new(CalendarWeek)
	err := c.getJSON(URL, &calendarWeek)
	if err != nil {
		return nil, err
	}

	return calendarWeek, nil
}
