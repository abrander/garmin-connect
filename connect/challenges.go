package main

import (
	"os"
	"strconv"
	"strings"

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

	challengesViewCmd := &cobra.Command{
		Use:  "view [id]",
		Run:  challengesView,
		Args: cobra.ExactArgs(1),
	}
	challengesCmd.AddCommand(challengesViewCmd)
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

func challengesView(_ *cobra.Command, args []string) {
	uuid := args[0]
	challenge, err := client.AdhocChallenge(uuid)
	bail(err)

	players := make([]string, len(challenge.Players), len(challenge.Players))
	for i, player := range challenge.Players {
		players[i] = player.FullName
	}

	t := NewTabular()
	t.AddValue("ID", challenge.UUID)
	t.AddValue("Start", challenge.Start.String())
	t.AddValue("End", challenge.End.String())
	t.AddValue("Description", challenge.Description)
	t.AddValue("Name", challenge.Name)
	t.AddValue("Rank", challenge.UserRanking)
	t.AddValue("Players", strings.Join(players, ", "))
	t.Output(os.Stdout)
}
