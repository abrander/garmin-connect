package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/abrander/garmin-connect"
)

func init() {
	badgesCmd := &cobra.Command{
		Use: "badges",
	}
	rootCmd.AddCommand(badgesCmd)

	badgesLeaderboardCmd := &cobra.Command{
		Use: "leaderboard",
		Run: badgesLeaderboard,
	}
	badgesCmd.AddCommand(badgesLeaderboardCmd)

	badgesEarnedCmd := &cobra.Command{
		Use:  "earned",
		Run:  badgesEarned,
		Args: cobra.RangeArgs(0, 1),
	}
	badgesCmd.AddCommand(badgesEarnedCmd)

	badgesCompareCmd := &cobra.Command{
		Use:  "compare [displayName]",
		Run:  badgesCompare,
		Args: cobra.ExactArgs(1),
	}
	badgesCmd.AddCommand(badgesCompareCmd)
}

func badgesLeaderboard(_ *cobra.Command, _ []string) {
	leaderboard, err := client.BadgeLeaderBoard()
	bail(err)

	t := NewTable()
	t.AddHeader("Display Name", "Name", "Level", "Points")
	for _, status := range leaderboard {
		t.AddRow(status.DisplayName, status.Fullname, strconv.Itoa(status.Level), strconv.Itoa(status.Point))
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
	t.AddHeader("Badge", "Points", "Date")
	for _, badge := range badges {
		t.AddRow(badge.Name, fmt.Sprintf("%d x%d", badge.Points, badge.EarnedNumber), badge.EarnedDate.String())
	}
	t.Output(os.Stdout)
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
			me = "✓"
			if e.meEarned > 1 {
				me += fmt.Sprintf(" %dx", e.meEarned)
			}
		}

		if e.other {
			other = "✓"
			if e.otherEarned > 1 {
				other += fmt.Sprintf(" %dx", e.otherEarned)
			}
		}
		t.AddRow(e.name, me, other, strconv.Itoa(e.points))
	}

	t.Output(os.Stdout)
}
