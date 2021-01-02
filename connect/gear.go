package main

import (
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

func init() {
	gearCmd := &cobra.Command{
		Use: "gear",
	}
	rootCmd.AddCommand(gearCmd)

	gearListCmd := &cobra.Command{
		Use:   "list [profile ID]",
		Short: "List Gear",
		Run:   gearList,
		Args:  cobra.RangeArgs(0, 1),
	}
	gearCmd.AddCommand(gearListCmd)
}

func gearList(_ *cobra.Command, args []string) {
	var profileID int64 = 0
	var err error
	if len(args) == 1 {
		profileID, err = strconv.ParseInt(args[0], 10, 64)
		bail(err)
	}
	gear, err := client.Gear(profileID)
	bail(err)

	t := NewTable()
	t.AddHeader("UUID", "Make", "Model", "Type", "Name", "Created Date")
	for _, g := range gear {
		t.AddRow(
			g.Uuid,
			g.GearMakeName,
			g.GearModelName,
			g.GearTypeName,
			g.DisplayName,
			g.CreateDate.Time,
		)
	}
	t.Output(os.Stdout)
}
