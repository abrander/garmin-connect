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
		Use:   "list",
		Short: "List ad-hoc challenges",
		Run:   challengesList,
		Args:  cobra.NoArgs,
	}
	challengesCmd.AddCommand(challengesListCmd)

	challengesListInvitesCmd := &cobra.Command{
		Use:   "invites",
		Short: "List ad-hoc challenge invites",
		Run:   challengesListInvites,
		Args:  cobra.NoArgs,
	}
	challengesListCmd.AddCommand(challengesListInvitesCmd)

	challengesAcceptCmd := &cobra.Command{
		Use:   "accept <invation ID>",
		Short: "Accept an ad-hoc challenge",
		Run:   challengesAccept,
		Args:  cobra.ExactArgs(1),
	}
	challengesCmd.AddCommand(challengesAcceptCmd)

	challengesDeclineCmd := &cobra.Command{
		Use:   "decline <invation ID>",
		Short: "Decline an ad-hoc challenge",
		Run:   challengesDecline,
		Args:  cobra.ExactArgs(1),
	}
	challengesCmd.AddCommand(challengesDeclineCmd)

	challengesListPreviousCmd := &cobra.Command{
		Use:   "previous",
		Short: "Show completed ad-hoc challenges",
		Run:   challengesListPrevious,
		Args:  cobra.NoArgs,
	}
	challengesListCmd.AddCommand(challengesListPreviousCmd)

	challengesViewCmd := &cobra.Command{
		Use:   "view <id>",
		Short: "View challenge details",
		Run:   challengesView,
		Args:  cobra.ExactArgs(1),
	}
	challengesCmd.AddCommand(challengesViewCmd)

	challengesLeaveCmd := &cobra.Command{
		Use:   "leave <challenge id>",
		Short: "Leave a challenge",
		Run:   challengesLeave,
		Args:  cobra.ExactArgs(1),
	}
	challengesCmd.AddCommand(challengesLeaveCmd)

	challengesRemoveCmd := &cobra.Command{
		Use:   "remove <challenge id> <user id>",
		Short: "Remove a user from a challenge",
		Run:   challengesRemove,
		Args:  cobra.ExactArgs(2),
	}
	challengesCmd.AddCommand(challengesRemoveCmd)
}

func challengesList(_ *cobra.Command, args []string) {
	challenges, err := client.AdhocChallenges()
	bail(err)

	t := NewTable()
	t.AddHeader("ID", "Start", "End", "Description", "Name", "Rank")
	for _, c := range challenges {
		t.AddRow(c.UUID, c.Start, c.End, c.Description, c.Name, c.UserRanking)
	}
	t.Output(os.Stdout)
}

func challengesListInvites(_ *cobra.Command, _ []string) {
	challenges, err := client.AdhocChallengeInvites()
	bail(err)

	t := NewTable()
	t.AddHeader("Invite ID", "Challenge ID", "Start", "End", "Description", "Name", "Rank")
	for _, c := range challenges {
		t.AddRow(c.InviteID, c.UUID, c.Start, c.End, c.Description, c.Name, c.UserRanking)
	}
	t.Output(os.Stdout)
}

func challengesAccept(_ *cobra.Command, args []string) {
	inviteID, err := strconv.Atoi(args[0])
	bail(err)

	err = client.AdhocChallengeInvitationRespond(inviteID, true)
	bail(err)
}

func challengesDecline(_ *cobra.Command, args []string) {
	inviteID, err := strconv.Atoi(args[0])
	bail(err)

	err = client.AdhocChallengeInvitationRespond(inviteID, false)
	bail(err)
}

func challengesListPrevious(_ *cobra.Command, args []string) {
	challenges, err := client.HistoricalAdhocChallenges()
	bail(err)

	t := NewTable()
	t.AddHeader("ID", "Start", "End", "Description", "Name", "Rank")
	for _, c := range challenges {
		t.AddRow(c.UUID, c.Start, c.End, c.Description, c.Name, c.UserRanking)
	}
	t.Output(os.Stdout)
}

func challengesLeave(_ *cobra.Command, args []string) {
	uuid := args[0]
	err := client.LeaveAdhocChallenge(uuid, 0)
	bail(err)
}

func challengesRemove(_ *cobra.Command, args []string) {
	uuid := args[0]

	profileID, err := strconv.ParseInt(args[1], 10, 64)
	bail(err)

	err = client.LeaveAdhocChallenge(uuid, profileID)
	bail(err)
}

func challengesView(_ *cobra.Command, args []string) {
	uuid := args[0]
	challenge, err := client.AdhocChallenge(uuid)
	bail(err)

	players := make([]string, len(challenge.Players))
	for i, player := range challenge.Players {
		players[i] = player.FullName + " [" + player.DisplayName + "]"
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
