package connect

import (
	"encoding/json"
	"fmt"
	"time"
)

type SportType struct {
	SportTypeId  int    `json:"sportTypeId"`
	SportTypeKey string `json:"sportTypeKey"`
}

type Author struct {
	UserProfilePk        int    `json:"userProfilePk"`
	DisplayName          string `json:"displayName"`
	FullName             string `json:"fullName"`
	ProfileImgNameLarge  string `json:"profileImgNameLarge"`
	ProfileImgNameMedium string `json:"profileImgNameMedium"`
	ProfileImgNameSmall  string `json:"profileImgNameSmall"`
	UserPro              bool   `json:"userPro"`
	VivokidUser          bool   `json:"vivokidUser"`
}

type EstimatedDistanceUnit struct {
	UnitId  int     `json:"unitId"`
	UnitKey string  `json:"unitKey"`
	Factor  float64 `json:"factor"`
}

type WorkoutSegment struct {
	SegmentOrder int           `json:"segmentOrder"`
	SportType    *SportType    `json:"sportType"`
	WorkoutSteps []WorkoutStep `json:"workoutSteps"`
}

type StepType struct {
	StepTypeId  int    `json:"stepTypeId"`
	StepTypeKey string `json:"stepTypeKey"`
}

type EndCondition struct {
	ConditionTypeId  int    `json:"conditionTypeId"`
	ConditionTypeKey string `json:"conditionTypeKey"`
	Displayable      bool   `json:"displayable"`
}

type TargetType struct {
	WorkoutTargetTypeId  int    `json:"workoutTargetTypeId"`
	WorkoutTargetTypeKey string `json:"workoutTargetTypeKey"`
}

type PreferredEndConditionUnit struct {
	UnitId  int     `json:"unitId"`
	UnitKey string  `json:"unitKey"`
	Factor  float64 `json:"factor"`
}

type WorkoutStep struct {
	Type                      string                     `json:"type"`
	StepId                    int                        `json:"stepId"`
	StepOrder                 int                        `json:"stepOrder"`
	StepType                  *StepType                  `json:"stepType"`
	ChildStepId               int                        `json:"childStepId"`
	Description               string                     `json:"description"`
	EndCondition              *EndCondition              `json:"endCondition"`
	EndConditionValue         float64                    `json:"endConditionValue"`
	PreferredEndConditionUnit *PreferredEndConditionUnit `json:"preferredEndConditionUnit,omitempty"`
	EndConditionCompare       bool                       `json:"endConditionCompare"`
	TargetType                *TargetType                `json:"targetType"`
	TargetValueOne            float64                    `json:"targetValueOne,omitempty"`
	TargetValueTwo            float64                    `json:"targetValueTwo,omitempty"`
	TargetValueUnit           string                     `json:"targetValueUnit,omitempty"`
	ZoneNumber                int                        `json:"zoneNumber"`
	// Various others..
}

// Workout describes a Garmin Connect workout entry
type Workout struct {
	WorkoutId                 int                    `json:"workoutId"`
	WorkoutName               string                 `json:"workoutName"`
	OwnerId                   int                    `json:"ownerId"`
	Description               *string                `json:"description,omitempty"`
	UpdateDate                *Time                  `json:"updateDate,omitempty"`
	CreatedDate               *Time                  `json:"createdDate,omitempty"`
	SportType                 *SportType             `json:"sportType,omitempty"`
	TrainingPlanId            *int                   `json:"trainingPlanId,omitempty"`
	Author                    *Author                `json:"author,omitempty"`
	EstimatedDurationInSecs   *int                   `json:"estimatedDurationInSecs,omitempty"`
	EstimatedDistanceInMeters *float64               `json:"estimatedDistanceInMeters,omitempty"`
	WorkoutSegments           []*WorkoutSegment      `json:"workoutSegments,omitempty"`
	EstimateType              *string                `json:"estimateType,omitempty"`
	EstimatedDistanceUnit     *EstimatedDistanceUnit `json:"estimatedDistanceUnit,omitempty"`
	Locale                    *string                `json:"locale,omitempty"`
	WorkoutProvider           *string                `json:"workoutProvider,omitempty"`
	UploadTimestamp           *Time                  `json:"uploadTimestamp,omitempty"`
	Consumer                  *string                `json:"consumer,omitempty"`
	ConsumerName              *string                `json:"consumerName,omitempty"`
	ConsumerImageUrl          *string                `json:"consumerImageURL,omitempty"`
	ConsumerWebsiteUrl        *string                `json:"consumerWebsiteURL,omitempty"`
	AtpPlanId                 *int                   `json:"atpPlanId,omitempty"`
	WorkoutNameI18nKey        *string                `json:"workoutNameI18nKey,omitempty"`
	DescriptionI18nKey        *string                `json:"descriptionI18nKey,omitempty"`
	AvgTrainingSpeed          *float64               `json:"avgTrainingSpeed,omitempty"`
	Shared                    *bool                  `json:"shared,omitempty"`
}

// workoutRequest is the bare minimum required to create a workout
type workoutRequest struct {
	SportType        *SportType        `json:"sportType"`
	WorkoutName      string            `json:"workoutName"`
	WorkoutSegments  []*WorkoutSegment `json:"workoutSegments"`
	AvgTrainingSpeed float64           `json:"avgTrainingSpeed"`
	Description      string            `json:"description,omitempty"`
}

func (w *Workout) MarshalJSON() ([]byte, error) {
	type Alias Workout

	var createdDate string
	if w.CreatedDate == nil || w.CreatedDate.IsZero() {
		createdDate = ""
	} else {
		createdDate = w.CreatedDate.Format("2006-01-02T15:04:05.0")
	}

	return json.Marshal(&struct {
		*Alias
		CreatedDate string `json:"createdDate"`
	}{
		Alias:       (*Alias)(w),
		CreatedDate: createdDate,
	})
}

func (c *Client) Workout() ([]Workout, error) {
	URL := "https://connect.garmin.com/modern/proxy/workout-service/workouts"
	var workout []Workout
	err := c.getJSON(URL, &workout)
	if err != nil {
		return nil, err
	}

	return workout, nil
}

func (c *Client) GetWorkout(workoutId int) (*Workout, error) {
	workout := new(Workout)
	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/workout-service/workout/%d", workoutId)
	err := c.getJSON(URL, &workout)
	if err != nil {
		return nil, err
	}

	return workout, nil
}

func (c *Client) UpdateWorkout(workout *Workout) error {
	workout.UpdateDate = nil

	workoutId := workout.WorkoutId
	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/workout-service/workout/%d", workoutId)
	err := c.writeWithMethodOverride("POST", URL, workout, 204, "PUT")
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) CreateWorkout(workout *Workout) (*Workout, error) {
	workoutRequest := &workoutRequest{
		SportType:        workout.SportType,
		WorkoutName:      workout.WorkoutName,
		WorkoutSegments:  workout.WorkoutSegments,
		AvgTrainingSpeed: *workout.AvgTrainingSpeed,
		Description:      *workout.Description,
	}

	fmt.Println(workoutRequest)
	workoutResponse := new(Workout)
	URL := "https://connect.garmin.com/modern/proxy/workout-service/workout"
	err := c.writeAndGetJSON("POST", URL, workoutRequest, 200, &workoutResponse)
	if err != nil {
		return nil, err
	}

	return workoutResponse, nil
}

type WorkoutSchedule struct {
	WorkoutScheduleId int     `json:"workoutScheduleId"`
	Workout           Workout `json:"workout"`
	CalendarDate      string  `json:"calendarDate"`
	CreatedDate       string  `json:"createdDate"`
	OwnerId           int     `json:"ownerId"`
}

type WorkoutSchedulePayload struct {
	Date *time.Time
}

func (s *WorkoutSchedulePayload) MarshalJSON() ([]byte, error) {
	type Alias WorkoutSchedulePayload

	var date string
	if s.Date == nil || s.Date.IsZero() {
		date = ""
	} else {
		date = s.Date.Format("2006-01-02")
	}

	return json.Marshal(&struct {
		*Alias
		Date string `json:"date"`
	}{
		Alias: (*Alias)(s),
		Date:  date,
	})
}

func (c *Client) ScheduleWorkout(workoutId int, date *time.Time) (*WorkoutSchedule, error) {
	workoutSchedule := new(WorkoutSchedule)
	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/workout-service/schedule/%d", workoutId)

	payload := &WorkoutSchedulePayload{Date: date}

	err := c.writeAndGetJSON("POST", URL, payload, 200, &workoutSchedule)
	if err != nil {
		return nil, err
	}

	return workoutSchedule, nil
}

type WorkoutScheduleSummary struct {
	ScheduleId int
	WorkoutId  int
	Title      string
	Date       string
}

func (c *Client) WorkoutSchedule(year, month int) ([]*WorkoutScheduleSummary, error) {
	calendarMonth, err := c.CalendarMonth(year, month)
	if err != nil {
		return nil, err
	}

	var workouts []*WorkoutScheduleSummary
	for _, activity := range calendarMonth.CalendarItems {
		if activity.ItemType == "workout" {
			workouts = append(workouts, &WorkoutScheduleSummary{
				ScheduleId: activity.ID,
				WorkoutId:  activity.WorkoutId,
				Title:      activity.Title,
				Date:       activity.Date.String(),
			})
		}
	}

	return workouts, nil
}

func (c *Client) DeleteScheduledWorkout(workoutId int) error {
	URL := fmt.Sprintf("https://connect.garmin.com/modern/proxy/workout-service/schedule/%d", workoutId)
	err := c.write("DELETE", URL, nil, 200)
	if err != nil {
		return err
	}

	return nil
}
