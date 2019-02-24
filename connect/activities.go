package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	activitiesCmd := &cobra.Command{
		Use: "activities",
	}
	rootCmd.AddCommand(activitiesCmd)

	activitiesListCmd := &cobra.Command{
		Use: "list",
		Run: activitiesList,
	}
	activitiesCmd.AddCommand(activitiesListCmd)
}

func activitiesList(_ *cobra.Command, args []string) {
	displayName := ""

	activities, err := client.Activities(displayName, 0, 100)
	bail(err)

	for _, a := range activities {
		fmt.Printf("%s %s %.01f\n", a.ActivityName, a.ActivityType.TypeKey, a.Distance)
	}
}
