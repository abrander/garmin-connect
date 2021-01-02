package main

import (
	"os"
	"sort"
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

	gearTypeListCmd := &cobra.Command{
		Use:   "type-list",
		Short: "List Gear Types",
		Run:   gearTypeList,
	}
	gearCmd.AddCommand(gearTypeListCmd)
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
	t.AddHeader("UUID", "Type", "Brand & Model", "Nickname", "Created Date", "Total Distance", "Activities")
	for _, g := range gear {

		gearStats, err := client.GearStats(g.Uuid)
		bail(err)

		//totalDistance, err := gearStats.TotalDistance.Int64()
		//bail(err)

		t.AddRow(
			g.Uuid,
			g.GearTypeName,
			g.CustomMakeModel,
			g.DisplayName,
			g.CreateDate.Time,
			strconv.FormatFloat(gearStats.TotalDistance, 'f', 2, 64),
			gearStats.TotalActivities,
		)
	}
	t.Output(os.Stdout)
}

func gearTypeList(_ *cobra.Command, args []string) {
	gearTypes, err := client.GearType()
	bail(err)

	t := NewTable()
	t.AddHeader("ID", "Name", "Created Date", "Update Date")
	sort.Slice(gearTypes, func(i, j int) bool {
		return gearTypes[i].TypeID < gearTypes[j].TypeID
	})

	for _, g := range gearTypes {
		t.AddRow(
			g.TypeID,
			g.TypeName,
			g.CreateDate,
			g.UpdateDate,
		)
	}
	t.Output(os.Stdout)
}
