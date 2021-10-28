package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	connect "github.com/abrander/garmin-connect"
)

var (
	exportFormat string
	offset       int
	count        int
)

func init() {
	activitiesCmd := &cobra.Command{
		Use: "activities",
	}
	rootCmd.AddCommand(activitiesCmd)

	activitiesListCmd := &cobra.Command{
		Use:   "list [display name]",
		Short: "List Activities",
		Run:   activitiesList,
		Args:  cobra.RangeArgs(0, 1),
	}
	activitiesListCmd.Flags().IntVarP(&offset, "offset", "o", 0, "Paginating index where the list starts from")
	activitiesListCmd.Flags().IntVarP(&count, "count", "c", 100, "Count of elements to return")
	activitiesCmd.AddCommand(activitiesListCmd)

	activitiesViewCmd := &cobra.Command{
		Use:   "view <activity id>",
		Short: "View details for an activity",
		Run:   activitiesView,
		Args:  cobra.ExactArgs(1),
	}
	activitiesCmd.AddCommand(activitiesViewCmd)

	activitiesViewWeatherCmd := &cobra.Command{
		Use:   "weather <activity id>",
		Short: "View weather for an activity",
		Run:   activitiesViewWeather,
		Args:  cobra.ExactArgs(1),
	}
	activitiesViewCmd.AddCommand(activitiesViewWeatherCmd)

	activitiesViewHRZonesCmd := &cobra.Command{
		Use:   "hrzones <activity id>",
		Short: "View hr zones for an activity",
		Run:   activitiesViewHRZones,
		Args:  cobra.ExactArgs(1),
	}
	activitiesViewCmd.AddCommand(activitiesViewHRZonesCmd)

	activitiesExportCmd := &cobra.Command{
		Use:   "export <activity id>",
		Short: "Export an activity to a file",
		Run:   activitiesExport,
		Args:  cobra.ExactArgs(1),
	}
	activitiesExportCmd.Flags().StringVarP(&exportFormat, "format", "f", "fit", "Format of export (fit, tcx, gpx, kml, csv)")
	activitiesCmd.AddCommand(activitiesExportCmd)

	activitiesImportCmd := &cobra.Command{
		Use:   "import <path>",
		Short: "Import an activity from a file",
		Run:   activitiesImport,
		Args:  cobra.ExactArgs(1),
	}
	activitiesCmd.AddCommand(activitiesImportCmd)

	activitiesDeleteCmd := &cobra.Command{
		Use:   "delete <activity id>",
		Short: "Delete an activity",
		Run:   activitiesDelete,
		Args:  cobra.ExactArgs(1),
	}
	activitiesCmd.AddCommand(activitiesDeleteCmd)

	activitiesRenameCmd := &cobra.Command{
		Use:   "rename <activity id> <new name>",
		Short: "Rename an activity",
		Run:   activitiesRename,
		Args:  cobra.ExactArgs(2),
	}
	activitiesCmd.AddCommand(activitiesRenameCmd)
}

func activitiesList(_ *cobra.Command, args []string) {
	displayName := ""

	if len(args) == 1 {
		displayName = args[0]
	}

	activities, err := client.Activities(displayName, offset, count)
	bail(err)

	t := NewTable()
	t.AddHeader("ID", "Date", "Name", "Type", "Distance", "Time", "Avg/Max HR", "Calories")
	for _, a := range activities {
		t.AddRow(
			a.ID,
			a.StartLocal.Time,
			a.ActivityName,
			a.ActivityType.TypeKey,
			a.Distance,
			a.StartLocal,
			fmt.Sprintf("%.0f/%.0f", a.AverageHeartRate, a.MaxHeartRate),
			a.Calories,
		)
	}
	t.Output(os.Stdout)
}

func activitiesView(_ *cobra.Command, args []string) {
	activityID, err := strconv.Atoi(args[0])
	bail(err)

	activity, err := client.Activity(activityID)
	bail(err)

	t := NewTabular()
	t.AddValue("ID", activity.ID)
	t.AddValue("Name", activity.ActivityName)
	t.Output(os.Stdout)
}

func activitiesViewWeather(_ *cobra.Command, args []string) {
	activityID, err := strconv.Atoi(args[0])
	bail(err)

	weather, err := client.ActivityWeather(activityID)
	bail(err)

	t := NewTabular()
	t.AddValueUnit("Temperature", weather.Temperature, "°F")
	t.AddValueUnit("Apparent Temperature", weather.ApparentTemperature, "°F")
	t.AddValueUnit("Dew Point", weather.DewPoint, "°F")
	t.AddValueUnit("Relative Humidity", weather.RelativeHumidity, "%")
	t.AddValueUnit("Wind Direction", weather.WindDirection, weather.WindDirectionCompassPoint)
	t.AddValueUnit("Wind Speed", weather.WindSpeed, "mph")
	t.AddValue("Latitude", weather.Latitude)
	t.AddValue("Longitude", weather.Longitude)
	t.Output(os.Stdout)
}

func activitiesViewHRZones(_ *cobra.Command, args []string) {
	activityID, err := strconv.Atoi(args[0])
	bail(err)

	zones, err := client.ActivityHrZones(activityID)
	bail(err)

	t := NewTabular()
	//for (zone in zones)
	for i := 0; i < len(zones)-1; i++ {
		t.AddValue(fmt.Sprintf("Zone %d (%3d-%3dbpm)", zones[i].ZoneNumber, zones[i].ZoneLowBoundary, zones[i+1].ZoneLowBoundary),
			zones[i].TimeInZone)
	}
	t.AddValue(fmt.Sprintf("Zone %d ( > %dbpm )", zones[len(zones)-1].ZoneNumber, zones[len(zones)-1].ZoneLowBoundary),
		zones[len(zones)-1].TimeInZone)

	t.Output(os.Stdout)
}

func activitiesExport(_ *cobra.Command, args []string) {
	format, err := connect.FormatFromExtension(exportFormat)
	bail(err)

	activityID, err := strconv.Atoi(args[0])
	bail(err)

	name := fmt.Sprintf("%d.%s", activityID, format.Extension())
	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	bail(err)

	err = client.ExportActivity(activityID, f, format)
	bail(err)
}

func activitiesImport(_ *cobra.Command, args []string) {
	filename := args[0]

	f, err := os.Open(filename)
	bail(err)

	format, err := connect.FormatFromFilename(filename)
	bail(err)

	id, err := client.ImportActivity(f, format)
	bail(err)

	fmt.Printf("Activity ID %d imported\n", id)
}

func activitiesDelete(_ *cobra.Command, args []string) {
	activityID, err := strconv.Atoi(args[0])
	bail(err)

	err = client.DeleteActivity(activityID)
	bail(err)
}

func activitiesRename(_ *cobra.Command, args []string) {
	activityID, err := strconv.Atoi(args[0])
	bail(err)

	newName := args[1]

	err = client.RenameActivity(activityID, newName)
	bail(err)
}
