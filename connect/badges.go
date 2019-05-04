package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	connect "github.com/abrander/garmin-connect"
)

const gotIt = "✓"

func init() {
	badgesCmd := &cobra.Command{
		Use: "badges",
	}
	rootCmd.AddCommand(badgesCmd)

	badgesLeaderboardCmd := &cobra.Command{
		Use:   "leaderboard",
		Short: "Show the current points leaderbaord among the authenticated users connections",
		Run:   badgesLeaderboard,
		Args:  cobra.NoArgs,
	}
	badgesCmd.AddCommand(badgesLeaderboardCmd)

	badgesEarnedCmd := &cobra.Command{
		Use:   "earned [display name]",
		Short: "Show the earned badges",
		Run:   badgesEarned,
		Args:  cobra.RangeArgs(0, 1),
	}
	badgesCmd.AddCommand(badgesEarnedCmd)

	badgesAvailableCmd := &cobra.Command{
		Use:   "available",
		Short: "Show badges not yet earned",
		Run:   badgesAvailable,
		Args:  cobra.NoArgs,
	}
	badgesCmd.AddCommand(badgesAvailableCmd)

	badgesViewCmd := &cobra.Command{
		Use:   "view <badge id>",
		Short: "Show details about a badge",
		Run:   badgesView,
		Args:  cobra.ExactArgs(1),
	}
	badgesCmd.AddCommand(badgesViewCmd)

	badgesCompareCmd := &cobra.Command{
		Use:   "compare <display name>",
		Short: "Compare the authenticated users badges with the badges of another user",
		Run:   badgesCompare,
		Args:  cobra.ExactArgs(1),
	}
	badgesCmd.AddCommand(badgesCompareCmd)
}

func badgesLeaderboard(_ *cobra.Command, _ []string) {
	leaderboard, err := client.BadgeLeaderBoard()
	bail(err)

	t := NewTable()
	t.AddHeader("Display Name", "Name", "Level", "Points")
	for _, status := range leaderboard {
		t.AddRow(status.DisplayName, status.Fullname, status.Level, status.Point)
	}
	t.Output(os.Stdout)
}

func badgesEarned(_ *cobra.Command, args []string) {
	var badges []connect.Badge

	if len(args) == 1 {
		displayName := args[0]
		// If we have a displayid to show, we abuse the compare call to read
		// badges earned by a connection.
		_, status, err := client.BadgeCompare(displayName)
		bail(err)

		badges = status.Badges
	} else {
		var err error
		badges, err = client.BadgesEarned()
		bail(err)
	}

	t := NewTable()
	t.AddHeader("ID", "Badge", "Points", "Date")
	for _, badge := range badges {
		p := fmt.Sprintf("%d", badge.Points)
		if badge.EarnedNumber > 1 {
			p = fmt.Sprintf("%d x%d", badge.Points, badge.EarnedNumber)
		}
		t.AddRow(badge.ID, badge.Name, p, badge.EarnedDate.String())
	}
	t.Output(os.Stdout)
}

func badgesAvailable(_ *cobra.Command, _ []string) {
	badges, err := client.BadgesAvailable()
	bail(err)

	t := NewTable()
	t.AddHeader("ID", "Key", "Name", "Points")
	for _, badge := range badges {
		t.AddRow(badge.ID, badge.Key, badge.Name, badge.Points)
	}
	t.Output(os.Stdout)
}

func badgesView(_ *cobra.Command, args []string) {
	badgeID, err := strconv.Atoi(args[0])
	bail(err)

	badge, err := client.BadgeDetail(badgeID)
	bail(err)

	t := NewTabular()
	t.AddValue("ID", badge.ID)
	t.AddValue("Key", badge.Key)
	t.AddValue("Name", badge.Name)
	t.AddValue("Points", badge.Points)
	t.AddValue("Earned", formatDate(badge.EarnedDate.Time))
	t.AddValueUnit("Earned", badge.EarnedNumber, "time(s)")
	t.AddValue("Available from", formatDate(badge.Start.Time))
	t.AddValue("Available to", formatDate(badge.End.Time))
	t.Output(os.Stdout)

	if len(badge.Connections) > 0 {
		fmt.Printf("\n  Connections with badge:\n")
		t := NewTable()
		t.AddHeader("Display Name", "Name", "Earned")
		for _, b := range badge.Connections {
			t.AddRow(b.DisplayName, b.FullName, b.EarnedDate.Time)
		}
		t.Output(os.Stdout)
	}

	if len(badge.RelatedBadges) > 0 {
		fmt.Printf("\n  Relates badges:\n")

		t := NewTable()
		t.AddHeader("ID", "Key", "Name", "Points", "Earned")
		for _, b := range badge.RelatedBadges {
			earned := ""
			if b.EarnedByMe {
				earned = gotIt
			}
			t.AddRow(b.ID, b.Key, b.Name, b.Points, earned)
		}
		t.Output(os.Stdout)
	}
}

func badgesCompare(_ *cobra.Command, args []string) {
	displayName := args[0]
	a, b, err := client.BadgeCompare(displayName)
	bail(err)

	t := NewTable()
	t.AddHeader("Badge", a.Fullname, b.Fullname, "Points")

	type status struct {
		name        string
		points      int
		me          bool
		meEarned    int
		other       bool
		otherEarned int
	}

	m := map[string]*status{}

	for _, badge := range a.Badges {
		s, found := m[badge.Key]
		if !found {
			s = &status{}
			m[badge.Key] = s
		}
		s.me = true
		s.meEarned = badge.EarnedNumber
		s.name = badge.Name
		s.points = badge.Points
	}

	for _, badge := range b.Badges {
		s, found := m[badge.Key]
		if !found {
			s = &status{}
			m[badge.Key] = s
		}
		s.other = true
		s.otherEarned = badge.EarnedNumber
		s.name = badge.Name
		s.points = badge.Points
	}

	for _, e := range m {
		var me string
		var other string
		if e.me {
			me = gotIt
			if e.meEarned > 1 {
				me += fmt.Sprintf(" %dx", e.meEarned)
			}
		}

		if e.other {
			other = gotIt
			if e.otherEarned > 1 {
				other += fmt.Sprintf(" %dx", e.otherEarned)
			}
		}
		t.AddRow(e.name, me, other, e.points)
	}

	t.Output(os.Stdout)
}
