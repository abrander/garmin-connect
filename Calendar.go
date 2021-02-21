package connect

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
	ElapsedDuration          int     `json:"elapsedDuration"`
	Strokes                  int     `json:"strokes"`
	MaxSpeed                 float64 `json:"maxSpeed"`
	ShareableEvent           bool    `json:"shareableEvent"`
	AutoCalcCalories         bool    `json:"autoCalcCalories"`
	ProtectedWorkoutSchedule bool    `json:"protectedWorkoutSchedule"`
	IsParent                 bool    `json:"isParent"`
}
