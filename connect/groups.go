package main

import (
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

func init() {
	groupsCmd := &cobra.Command{
		Use: "groups",
	}
	rootCmd.AddCommand(groupsCmd)

	groupsListCmd := &cobra.Command{
		Use:   "list [display name]",
		Short: "List all groups",
		Run:   groupsList,
		Args:  cobra.RangeArgs(0, 1),
	}
	groupsCmd.AddCommand(groupsListCmd)

	groupsViewCmd := &cobra.Command{
		Use:   "view <group name>",
		Short: "View group details",
		Run:   groupsView,
		Args:  cobra.ExactArgs(1),
	}
	groupsCmd.AddCommand(groupsViewCmd)

	groupsSearchCmd := &cobra.Command{
		Use:   "search <keyword>",
		Short: "Search for a group",
		Run:   groupsSearch,
		Args:  cobra.ExactArgs(1),
	}
	groupsCmd.AddCommand(groupsSearchCmd)

	groupsJoinCmd := &cobra.Command{
		Use:   "join <profile id> <group id>",
		Short: "Join a group",
		Run:   groupsJoin,
		Args:  cobra.ExactArgs(2),
	}
	groupsCmd.AddCommand(groupsJoinCmd)

	groupsLeaveCmd := &cobra.Command{
		Use:   "leave <profile id> <group id>",
		Short: "Leave a group",
		Run:   groupsLeave,
		Args:  cobra.ExactArgs(2),
	}
	groupsCmd.AddCommand(groupsLeaveCmd)
}

func groupsList(_ *cobra.Command, args []string) {
	displayName := ""
	if len(args) == 1 {
		displayName = args[0]
	}

	groups, err := client.Groups(displayName)
	bail(err)

	t := NewTable()
	t.AddHeader("ID", "Name", "Description", "Profile Image")
	for _, g := range groups {
		t.AddRow(strconv.Itoa(g.ID), g.Name, g.Description, g.ProfileImageURLLarge)
	}
	t.Output(os.Stdout)
}

func groupsSearch(_ *cobra.Command, args []string) {
	keyword := args[0]
	groups, err := client.SearchGroups(keyword)
	bail(err)

	t := NewTable()
	t.AddHeader("ID", "Name", "Description", "Profile Image")
	for _, g := range groups {
		t.AddRow(strconv.Itoa(g.ID), g.Name, g.Description, g.ProfileImageURLLarge)
	}
	t.Output(os.Stdout)
}

func groupsView(_ *cobra.Command, args []string) {
	id, err := strconv.Atoi(args[0])
	bail(err)

	group, err := client.Group(id)
	bail(err)

	t := NewTabular()
	t.AddValue("ID", group.ID)
	t.AddValue("Name", group.Name)
	t.AddValue("Description", group.Description)
	t.AddValue("OwnerID", group.OwnerID)
	t.AddValue("ProfileImageURLLarge", group.ProfileImageURLLarge)
	t.AddValue("ProfileImageURLMedium", group.ProfileImageURLMedium)
	t.AddValue("ProfileImageURLSmall", group.ProfileImageURLSmall)
	t.AddValue("Visibility", group.Visibility)
	t.AddValue("Privacy", group.Privacy)
	t.AddValue("Location", group.Location)
	t.AddValue("WebsiteURL", group.WebsiteURL)
	t.AddValue("FacebookURL", group.FacebookURL)
	t.AddValue("TwitterURL", group.TwitterURL)
	//	t.AddValue("PrimaryActivities", group.PrimaryActivities)
	t.AddValue("OtherPrimaryActivity", group.OtherPrimaryActivity)
	//	t.AddValue("LeaderboardTypes", group.LeaderboardTypes)
	//	t.AddValue("FeatureTypes", group.FeatureTypes)
	t.AddValue("CorporateWellness", group.CorporateWellness)
	//	t.AddValue("ActivityFeedTypes", group.ActivityFeedTypes)
	t.Output(os.Stdout)
}

func groupsJoin(_ *cobra.Command, args []string) {
	profileID, err := strconv.Atoi(args[0])
	bail(err)

	groupID, err := strconv.Atoi(args[1])
	bail(err)

	err = client.JoinGroup(profileID, groupID)
	bail(err)
}

func groupsLeave(_ *cobra.Command, args []string) {
	profileID, err := strconv.Atoi(args[0])
	bail(err)

	groupID, err := strconv.Atoi(args[1])
	bail(err)

	err = client.LeaveGroup(profileID, groupID)
	bail(err)
}
