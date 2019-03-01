package main

import (
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

func init() {
	challengesCmd := &cobra.Command{
		Use: "challenges",
	}
	rootCmd.AddCommand(challengesCmd)

	challengesListCmd := &cobra.Command{
		Use:  "list",
		Run:  challengesList,
		Args: cobra.ExactArgs(0),
	}
	challengesCmd.AddCommand(challengesListCmd)
}

func challengesList(_ *cobra.Command, args []string) {
	challenges, err := client.AdhocChallenges()
	bail(err)

	t := NewTable()
	t.AddHeader("ID", "Start", "End", "Description", "Name", "Rank")
	for _, c := range challenges {
		t.AddRow(c.UUID, c.Start.String(), c.End.String(), c.Description, c.Name, strconv.Itoa(c.UserRanking))
	}
	t.Output(os.Stdout)
}
