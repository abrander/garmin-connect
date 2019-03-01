package main

import (
	"os"
	"strconv"
	"time"

	"github.com/abrander/garmin-connect"
	"github.com/spf13/cobra"
)

func init() {
	weightCmd := &cobra.Command{
		Use: "weight",
	}
	rootCmd.AddCommand(weightCmd)

	weightLatestCmd := &cobra.Command{
		Use: "latest",
		Run: weightLatest,
	}
	weightCmd.AddCommand(weightLatestCmd)

	weightAddCmd := &cobra.Command{
		Use:  "add [yyyy-mm-dd] [weight in grams]",
		Run:  weightAdd,
		Args: cobra.ExactArgs(2),
	}
	weightCmd.AddCommand(weightAddCmd)
}

func weightLatest(_ *cobra.Command, _ []string) {
	weightin, err := client.LatestWeight(time.Now())
	bail(err)

	t := NewTabular()
	t.AddValue("Date", weightin.Date.String())
	t.AddValueUnit("Weight", weightin.Weight/1000.0, "kg")
	t.AddValueUnit("BMI", weightin.BMI, "kg/m2")
	t.AddValueUnit("Fat", weightin.BodyFatPercentage, "%")
	t.AddValueUnit("Water", weightin.BodyWater, "%")
	t.AddValueUnit("Bone Mass", float64(weightin.BoneMass)/1000.0, "kg")
	t.AddValueUnit("Muscle Mass", float64(weightin.MuscleMass)/1000.0, "kg")
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
