package main

import (
	"os"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	infoCmd := &cobra.Command{
		Use:  "info",
		Run:  info,
		Args: cobra.RangeArgs(0, 1),
	}
	rootCmd.AddCommand(infoCmd)
}

func info(_ *cobra.Command, args []string) {
	displayName := ""
	if len(args) == 1 {
		displayName = args[0]
	}

	t := NewTabular()

	socialProfile, err := client.SocialProfile(displayName)
	bail(err)
	t.AddValue("ID", socialProfile.ID)
	t.AddValue("Display Name", socialProfile.DisplayName)
	t.AddValue("Name", socialProfile.Fullname)
	t.AddValue("Level", socialProfile.UserLevel)
	t.AddValue("Points", socialProfile.UserPoint)
	t.AddValue("Profile Image", socialProfile.ProfileImageURLLarge)

	info, err := client.PersonalInformation(socialProfile.DisplayName)
	bail(err)

	t.AddValue("", "")
	t.AddValue("Gender", info.UserInfo.Gender)
	t.AddValueUnit("Age", info.UserInfo.Age, "years")
	t.AddValueUnit("Height", info.BiometricProfile.Height, "cm")
	t.AddValueUnit("Weight", info.BiometricProfile.Weight/1000.0, "kg")
	t.AddValueUnit("Vo² Max", info.BiometricProfile.VO2Max, "mL/kg/min")
	t.AddValueUnit("Vo² Max (cycling)", info.BiometricProfile.VO2MaxCycling, "mL/kg/min")

	life, err := client.LifetimeActivities(socialProfile.DisplayName)
	bail(err)

	t.AddValue("", "")
	t.AddValue("Activities", life.Activities)
	t.AddValueUnit("Distance", life.Distance/1000.0, "km")
	t.AddValueUnit("Time", (time.Duration(life.Duration) * time.Second).Round(time.Second).String(), "hms")
	t.AddValueUnit("Calories", life.Calories/4.184, "Kcal")
	t.AddValueUnit("Elev Gain", life.ElevationGain, "m")

	totals, err := client.LifetimeTotals(socialProfile.DisplayName)
	bail(err)

	t.AddValue("", "")
	t.AddValueUnit("Steps", totals.Steps, "steps")
	t.AddValueUnit("Distance", totals.Distance/1000.0, "km")
	t.AddValueUnit("Daily Goal Met", totals.GoalsMetInDays, "days")
	t.AddValueUnit("Active Days", totals.ActiveDays, "days")
	t.AddValueUnit("Average Steps", totals.Steps/totals.ActiveDays, "steps")
	t.AddValueUnit("Calories", totals.Calories, "kCal")

	lastUsed, err := client.LastUsed(socialProfile.DisplayName)
	bail(err)

	t.AddValue("", "")
	t.AddValue("Device ID", lastUsed.DeviceID)
	t.AddValue("Device", lastUsed.DeviceName)
	t.AddValue("Time", lastUsed.DeviceUploadTime.String())
	t.AddValue("Ago", time.Since(lastUsed.DeviceUploadTime.Time).Round(time.Second).String())
	t.AddValue("Image", lastUsed.ImageURL)
	t.Output(os.Stdout)
}
