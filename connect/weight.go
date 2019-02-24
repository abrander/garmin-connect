package main

import (
	"os"
	"time"

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
}

func weightLatest(_ *cobra.Command, _ []string) {
	weightin, err := client.LatestWeight(time.Now())
	bail(err)

	t := NewTabular()
	t.AddValue("Date", weightin.Date.String())
	t.AddValueUnit("Weight", weightin.Weight/1000.0, "kg")
	t.AddValueUnit("BMI", weightin.BMI, "kg/m2")
	t.AddValueUnit("Fat", weightin.BodyFatPercentage, "%")
	t.AddValueUnit("Water", weightin.BodyWater, "kg")
	t.AddValueUnit("Bone Mass", float64(weightin.BoneMass)/1000.0, "kg")
	t.AddValueUnit("Muscle Mass", float64(weightin.MuscleMass)/1000.0, "kg")
	t.Output(os.Stdout)
}
