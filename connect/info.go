package main

import (
	"os"
	"time"

	connect "github.com/abrander/garmin-connect"
	"github.com/spf13/cobra"
)

func init() {
	infoCmd := &cobra.Command{
		Use:   "info [display name]",
		Short: "Show various information and statistics about a Connect User",
		Run:   info,
		Args:  cobra.RangeArgs(0, 1),
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
	if err == connect.ErrNotFound {
		bail(err)
	}

	if err == nil {
		displayName = socialProfile.DisplayName
	} else {
		socialProfile, err = client.PublicSocialProfile(displayName)
		bail(err)

		displayName = socialProfile.DisplayName
	}

	t.AddValue("ID", socialProfile.ID)
	t.AddValue("Profile ID", socialProfile.ProfileID)
	t.AddValue("Display Name", socialProfile.DisplayName)
	t.AddValue("Name", socialProfile.Fullname)
	t.AddValue("Level", socialProfile.UserLevel)
	t.AddValue("Points", socialProfile.UserPoint)
	t.AddValue("Profile Image", socialProfile.ProfileImageURLLarge)

	info, err := client.PersonalInformation(displayName)
	if err == nil {
		t.AddValue("", "")
		t.AddValue("Gender", info.UserInfo.Gender)
		t.AddValueUnit("Age", info.UserInfo.Age, "years")
		t.AddValueUnit("Height", nzf(info.BiometricProfile.Height), "cm")
		t.AddValueUnit("Weight", nzf(info.BiometricProfile.Weight/1000.0), "kg")
		t.AddValueUnit("Vo² Max", nzf(info.BiometricProfile.VO2Max), "mL/kg/min")
		t.AddValueUnit("Vo² Max (cycling)", nzf(info.BiometricProfile.VO2MaxCycling), "mL/kg/min")
	}

	life, err := client.LifetimeActivities(displayName)
	if err == nil {
		t.AddValue("", "")
		t.AddValue("Activities", life.Activities)
		t.AddValueUnit("Distance", life.Distance/1000.0, "km")
		t.AddValueUnit("Time", (time.Duration(life.Duration) * time.Second).Round(time.Second).String(), "hms")
		t.AddValueUnit("Calories", life.Calories/4.184, "Kcal")
		t.AddValueUnit("Elev Gain", life.ElevationGain, "m")
	}

	totals, err := client.LifetimeTotals(displayName)
	if err == nil {
		t.AddValue("", "")
		t.AddValueUnit("Steps", totals.Steps, "steps")
		t.AddValueUnit("Distance", totals.Distance/1000.0, "km")
		t.AddValueUnit("Daily Goal Met", totals.GoalsMetInDays, "days")
		t.AddValueUnit("Active Days", totals.ActiveDays, "days")
		if totals.ActiveDays > 0 {
			t.AddValueUnit("Average Steps", totals.Steps/totals.ActiveDays, "steps")
		}
		t.AddValueUnit("Calories", totals.Calories, "kCal")
	}

	lastUsed, err := client.LastUsed(displayName)
	if err == nil {
		t.AddValue("", "")
		t.AddValue("Device ID", lastUsed.DeviceID)
		t.AddValue("Device", lastUsed.DeviceName)
		t.AddValue("Time", lastUsed.DeviceUploadTime.String())
		t.AddValue("Ago", time.Since(lastUsed.DeviceUploadTime.Time).Round(time.Second).String())
		t.AddValue("Image", lastUsed.ImageURL)
	}

	t.Output(os.Stdout)
}
