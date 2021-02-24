package main

import (
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

func init() {
	calendarCmd := &cobra.Command{
		Use: "calendar",
	}
	rootCmd.AddCommand(calendarCmd)

	calendarYearCmd := &cobra.Command{
		Use:   "year <year>",
		Short: "List active days in the year",
		Run:   calendarYear,
		Args:  cobra.RangeArgs(1, 1),
	}
	calendarCmd.AddCommand(calendarYearCmd)

	calendarMonthCmd := &cobra.Command{
		Use:   "month <year> <month>",
		Short: "List active days in the month",
		Run:   calendarMonth,
		Args:  cobra.RangeArgs(2, 2),
	}
	calendarCmd.AddCommand(calendarMonthCmd)

	calendarWeekCmd := &cobra.Command{
		Use:   "week <year> <month> <day>",
		Short: "List active days in the week",
		Run:   calendarWeek,
		Args:  cobra.RangeArgs(3, 3),
	}
	calendarCmd.AddCommand(calendarWeekCmd)

}

func calendarYear(_ *cobra.Command, args []string) {
	year, err := strconv.ParseInt(args[0], 10, 32)
	bail(err)

	calendar, err := client.CalendarYear(int(year))
	bail(err)

	t := NewTable()
	t.AddHeader("ActivityType ID", "Number of Activities", "Total Distance", "Total Duration", "Total Calories")
	for _, summary := range calendar.YearSummaries {
		t.AddRow(
			summary.ActivityTypeID,
			summary.NumberOfActivities,
			summary.TotalDistance,
			summary.TotalDuration,
			summary.TotalCalories,
		)
	}
	t.Output(os.Stdout)
}

func calendarMonth(_ *cobra.Command, args []string) {
	year, err := strconv.ParseInt(args[0], 10, 32)
	bail(err)

	month, err := strconv.ParseInt(args[1], 10, 32)
	bail(err)

	calendar, err := client.CalendarMonth(int(year), int(month))
	bail(err)

	t := NewTable()
	t.AddHeader("ID", "Date", "Name", "Distance", "Time", "Calories")
	for _, item := range calendar.CalendarItems {
		t.AddRow(
			item.ID,
			item.Date,
			item.Title,
			item.Distance,
			item.ElapsedDuration,
			item.Calories,
		)
	}
	t.Output(os.Stdout)
}

func calendarWeek(_ *cobra.Command, args []string) {
	year, err := strconv.ParseInt(args[0], 10, 32)
	bail(err)

	month, err := strconv.ParseInt(args[1], 10, 32)
	bail(err)

	week, err := strconv.ParseInt(args[2], 10, 32)
	bail(err)

	calendar, err := client.CalendarWeek(int(year), int(month), int(week))
	bail(err)

	t := NewTable()
	t.AddHeader("ID", "Date", "Name", "Distance", "Time", "Calories")
	for _, item := range calendar.CalendarItems {
		t.AddRow(
			item.ID,
			item.Date,
			item.Title,
			item.Distance,
			item.ElapsedDuration,
			item.Calories,
		)
	}
	t.Output(os.Stdout)
}
