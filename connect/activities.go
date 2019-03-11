package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/abrander/garmin-connect"
)

var (
	exportFormat string
)

func init() {
	activitiesCmd := &cobra.Command{
		Use: "activities",
	}
	rootCmd.AddCommand(activitiesCmd)

	activitiesListCmd := &cobra.Command{
		Use:   "list [display name]",
		Short: "List Actvities",
		Run:   activitiesList,
		Args:  cobra.RangeArgs(0, 1),
	}
	activitiesCmd.AddCommand(activitiesListCmd)

	activitiesViewCmd := &cobra.Command{
		Use:   "view <activity id>",
		Short: "View details for an activity",
		Run:   activitiesView,
		Args:  cobra.ExactArgs(1),
	}
	activitiesCmd.AddCommand(activitiesViewCmd)

	activitiesExportCmd := &cobra.Command{
		Use:   "export <activity id>",
		Short: "Export an activity to a file",
		Run:   activitiesExport,
		Args:  cobra.ExactArgs(1),
	}
	activitiesExportCmd.Flags().StringVarP(&exportFormat, "format", "f", "fit", "Format of export (fit, tcx, gpx, kml, csv)")
	activitiesCmd.AddCommand(activitiesExportCmd)

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

	activities, err := client.Activities(displayName, 0, 100)
	bail(err)

	t := NewTable()
	t.AddHeader("ID", "Date", "Name", "Type", "Distance", "Time", "Avg/Max HR", "Calories")
	for _, a := range activities {
		t.AddRow(
			strconv.Itoa(a.ID),
			formatDate(a.StartLocal.Time),
			a.ActivityName,
			a.ActivityType.TypeKey,
			fmt.Sprintf("%.0f", a.Distance),
			formatTime(a.StartLocal.Time),
			fmt.Sprintf("%.0f/%.0f", a.AverageHeartRate, a.MaxHeartRate),
			fmt.Sprintf("%.0f", a.Calories),
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

func activitiesExport(_ *cobra.Command, args []string) {
	formatTable := map[string]int{
		"fit": connect.FormatFIT,
		"tcx": connect.FormatTCX,
		"gpx": connect.FormatGPX,
		"kml": connect.FormatKML,
		"csv": connect.FormatCSV,
	}

	filenameTable := map[string]string{
		"fit": "%d.zip",
		"tcx": "%d.tcx",
		"gpx": "%d.gpx",
		"kml": "%d.kml",
		"csv": "%d.csv",
	}

	format, found := formatTable[exportFormat]
	if !found {
		bail(errors.New(exportFormat))
	}

	activityID, err := strconv.Atoi(args[0])
	bail(err)

	name := fmt.Sprintf(filenameTable[exportFormat], activityID)
	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	bail(err)

	err = client.ExportActivity(activityID, f, format)
	bail(err)
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
