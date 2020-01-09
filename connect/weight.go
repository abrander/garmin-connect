package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	connect "github.com/abrander/garmin-connect"
	"github.com/spf13/cobra"
)

func init() {
	weightCmd := &cobra.Command{
		Use: "weight",
	}
	rootCmd.AddCommand(weightCmd)

	weightLatestCmd := &cobra.Command{
		Use:   "latest",
		Short: "Show the latest weight-in",
		Run:   weightLatest,
		Args:  cobra.NoArgs,
	}
	weightCmd.AddCommand(weightLatestCmd)

	weightLatestWeekCmd := &cobra.Command{
		Use:   "week",
		Short: "Show average weight for the latest week",
		Run:   weightLatestWeek,
		Args:  cobra.NoArgs,
	}
	weightLatestCmd.AddCommand(weightLatestWeekCmd)

	weightAddCmd := &cobra.Command{
		Use:   "add <yyyy-mm-dd> <weight in grams>",
		Short: "Add a simple weight for a specific date",
		Run:   weightAdd,
		Args:  cobra.ExactArgs(2),
	}
	weightCmd.AddCommand(weightAddCmd)

	weightDeleteCmd := &cobra.Command{
		Use:   "delete <yyyy-mm-dd]>",
		Short: "Delete a weight-in",
		Run:   weightDelete,
		Args:  cobra.ExactArgs(1),
	}
	weightCmd.AddCommand(weightDeleteCmd)

	weightDateCmd := &cobra.Command{
		Use:   "date [yyyy-mm-dd]",
		Short: "Show weight for a specific date",
		Run:   weightDate,
		Args:  cobra.ExactArgs(1),
	}
	weightCmd.AddCommand(weightDateCmd)

	weightRangeCmd := &cobra.Command{
		Use:   "range [yyyy-mm-dd] [yyyy-mm-dd]",
		Short: "Show weight for a date range",
		Run:   weightRange,
		Args:  cobra.ExactArgs(2),
	}
	weightCmd.AddCommand(weightRangeCmd)

	weightGoalCmd := &cobra.Command{
		Use:   "goal [displayName]",
		Short: "Show weight goal",
		Run:   weightGoal,
		Args:  cobra.RangeArgs(0, 1),
	}
	weightCmd.AddCommand(weightGoalCmd)

	weightGoalSetCmd := &cobra.Command{
		Use:   "set [goal in gram]",
		Short: "Set weight goal",
		Run:   weightGoalSet,
		Args:  cobra.ExactArgs(1),
	}
	weightGoalCmd.AddCommand(weightGoalSetCmd)
}

func weightLatest(_ *cobra.Command, _ []string) {
	weightin, err := client.LatestWeight(time.Now())
	bail(err)

	t := NewTabular()
	t.AddValue("Date", weightin.Date.String())
	t.AddValueUnit("Weight", weightin.Weight/1000.0, "kg")
	t.AddValueUnit("BMI", weightin.BMI, "kg/m2")
	t.AddValueUnit("Fat", weightin.BodyFatPercentage, "%")
	t.AddValueUnit("Fat Mass", (weightin.Weight*weightin.BodyFatPercentage)/100000.0, "kg")
	t.AddValueUnit("Water", weightin.BodyWater, "%")
	t.AddValueUnit("Bone Mass", float64(weightin.BoneMass)/1000.0, "kg")
	t.AddValueUnit("Muscle Mass", float64(weightin.MuscleMass)/1000.0, "kg")
	t.Output(os.Stdout)
}

func weightLatestWeek(_ *cobra.Command, _ []string) {
	now := time.Now()
	from := time.Now().Add(-24 * 6 * time.Hour)

	average, _, err := client.Weightins(from, now)
	bail(err)

	t := NewTabular()
	t.AddValue("Average from", formatDate(from))
	t.AddValueUnit("Weight", average.Weight/1000.0, "kg")
	t.AddValueUnit("BMI", average.BMI, "kg/m2")
	t.AddValueUnit("Fat", average.BodyFatPercentage, "%")
	t.AddValueUnit("Fat Mass", (average.Weight*average.BodyFatPercentage)/100000.0, "kg")
	t.AddValueUnit("Water", average.BodyWater, "%")
	t.AddValueUnit("Bone Mass", float64(average.BoneMass)/1000.0, "kg")
	t.AddValueUnit("Muscle Mass", float64(average.MuscleMass)/1000.0, "kg")
	t.Output(os.Stdout)
}

func weightAdd(_ *cobra.Command, args []string) {
	date, err := connect.ParseDate(args[0])
	bail(err)

	weight, err := strconv.Atoi(args[1])
	bail(err)

	err = client.AddUserWeight(date.Time(), float64(weight))
	bail(err)
}

func weightDelete(_ *cobra.Command, args []string) {
	date, err := connect.ParseDate(args[0])
	bail(err)

	err = client.DeleteWeightin(date.Time())
	bail(err)
}

func weightDate(_ *cobra.Command, args []string) {
	date, err := connect.ParseDate(args[0])
	bail(err)

	tim, weight, err := client.WeightByDate(date.Time())
	bail(err)

	zero := time.Time{}
	if tim.Time == zero {
		fmt.Printf("No weight ins on this date\n")
		os.Exit(1)
	}

	t := NewTabular()
	t.AddValue("Time", tim.String())
	t.AddValueUnit("Weight", weight/1000.0, "kg")
	t.Output(os.Stdout)
}

func weightRange(_ *cobra.Command, args []string) {
	from, err := connect.ParseDate(args[0])
	bail(err)

	to, err := connect.ParseDate(args[1])
	bail(err)

	average, weightins, err := client.Weightins(from.Time(), to.Time())
	bail(err)

	t := NewTabular()

	t.AddValueUnit("Weight", average.Weight/1000.0, "kg")
	t.AddValueUnit("BMI", average.BMI, "kg/m2")
	t.AddValueUnit("Fat", average.BodyFatPercentage, "%")
	t.AddValueUnit("Fat Mass", average.Weight*average.BodyFatPercentage/100000.0, "kg")
	t.AddValueUnit("Water", average.BodyWater, "%")
	t.AddValueUnit("Bone Mass", float64(average.BoneMass)/1000.0, "kg")
	t.AddValueUnit("Muscle Mass", float64(average.MuscleMass)/1000.0, "kg")
	fmt.Fprintf(os.Stdout, "        \033[1mAverage\033[0m\n")
	t.Output(os.Stdout)

	t2 := NewTable()
	t2.AddHeader("Date", "Weight", "BMI", "Fat%", "Fat", "Water%", "Bone Mass", "Muscle Mass")
	for _, weightin := range weightins {
		if weightin.Weight < 1.0 {
			continue
		}

		t2.AddRow(
			weightin.Date,
			weightin.Weight/1000.0,
			nzf(weightin.BMI),
			nzf(weightin.BodyFatPercentage),
			nzf(weightin.Weight*weightin.BodyFatPercentage/100000.0),
			nzf(weightin.BodyWater),
			nzf(float64(weightin.BoneMass)/1000.0),
			nzf(float64(weightin.MuscleMass)/1000.0),
		)
	}
	fmt.Fprintf(os.Stdout, "\n")
	t2.Output(os.Stdout)
}

func weightGoal(_ *cobra.Command, args []string) {
	displayName := ""

	if len(args) > 0 {
		displayName = args[0]
	}

	goal, err := client.WeightGoal(displayName)
	bail(err)

	t := NewTabular()
	t.AddValue("ID", goal.ID)
	t.AddValue("Created", goal.Created)
	t.AddValueUnit("Target", float64(goal.Value)/1000.0, "kg")
	t.Output(os.Stdout)
}

func weightGoalSet(_ *cobra.Command, args []string) {
	goal, err := strconv.Atoi(args[0])
	bail(err)

	err = client.SetWeightGoal(goal)
	bail(err)
}
