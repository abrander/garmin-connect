package main

import (
	"os"

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

	//UserInfo:{Gender:FEMALE Email: Locale:da TimeZone:Europe/Paris Age:37} BiometricProfile:{UserID:60610314 Height:180 Weight:85569.99969482422 VO2Max:34 VO2MaxCycling:10}}

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
	t.Output(os.Stdout)
}
