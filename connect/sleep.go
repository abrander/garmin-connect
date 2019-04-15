package main

import (
	"fmt"
	"os"

	connect "github.com/abrander/garmin-connect"
	"github.com/spf13/cobra"
)

func init() {
	sleepCmd := &cobra.Command{
		Use: "sleep",
	}
	rootCmd.AddCommand(sleepCmd)

	sleepSummaryCmd := &cobra.Command{
		Use:   "summary <date> [displayName]",
		Short: "Show sleep summary for date",
		Run:   sleepSummary,
		Args:  cobra.RangeArgs(1, 2),
	}
	sleepCmd.AddCommand(sleepSummaryCmd)
}

func sleepSummary(_ *cobra.Command, args []string) {
	date, err := connect.ParseDate(args[0])
	bail(err)

	displayName := ""

	if len(args) > 1 {
		displayName = args[1]
	}

	summary, _, levels, err := client.SleepData(displayName, date.Time())
	bail(err)

	t := NewTabular()
	t.AddValue("Start", summary.StartGMT)
	t.AddValue("End", summary.EndGMT)
	t.AddValue("Sleep", hoursAndMinutes(summary.Sleep))
	t.AddValue("Nap", hoursAndMinutes(summary.Nap))
	t.AddValue("Unmeasurable", hoursAndMinutes(summary.Unmeasurable))
	t.AddValue("Deep", hoursAndMinutes(summary.Deep))
	t.AddValue("Light", hoursAndMinutes(summary.Light))
	t.AddValue("REM", hoursAndMinutes(summary.REM))
	t.AddValue("Awake", hoursAndMinutes(summary.Awake))
	t.AddValue("Confirmed", summary.Confirmed)
	t.AddValue("Confirmation Type", summary.Confirmation)
	t.AddValue("REM Data", summary.REMData)
	t.Output(os.Stdout)

	fmt.Fprintf(os.Stdout, "\n")

	t2 := NewTable()
	t2.AddHeader("Start", "End", "State", "Duration")
	for _, l := range levels {
		t2.AddRow(l.Start, l.End, l.State, hoursAndMinutes(l.End.Sub(l.Start.Time)))
	}
	t2.Output(os.Stdout)
}
