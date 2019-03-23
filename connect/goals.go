package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

func init() {
	goalsCmd := &cobra.Command{
		Use: "goals",
	}
	rootCmd.AddCommand(goalsCmd)

	goalsListCmd := &cobra.Command{
		Use:   "list [display name]",
		Short: "List all goals",
		Run:   goalsList,
		Args:  cobra.RangeArgs(0, 1),
	}
	goalsCmd.AddCommand(goalsListCmd)

	goalsDeleteCmd := &cobra.Command{
		Use:   "delete <goal id>",
		Short: "Delete a goal",
		Run:   goalsDelete,
		Args:  cobra.ExactArgs(1),
	}
	goalsCmd.AddCommand(goalsDeleteCmd)
}

func goalsList(_ *cobra.Command, args []string) {
	displayName := ""
	if len(args) == 1 {
		displayName = args[0]
	}

	t := NewTable()
	t.AddHeader("ID", "Profile", "Category", "Type", "Start", "End", "Created", "Value")
	for typ := 0; typ <= 9; typ++ {
		goals, err := client.Goals(displayName, typ)
		bail(err)

		for _, g := range goals {
			t.AddRow(
				fmt.Sprintf("%d", g.ID),
				fmt.Sprintf("%d", g.ProfileID),
				strconv.Itoa(g.GoalCategory),
				g.GoalType.String(),
				g.Start.String(),
				g.End.String(),
				g.Created.String(),
				strconv.Itoa(g.Value),
			)
		}
	}
	t.Output(os.Stdout)
}

func goalsDelete(_ *cobra.Command, args []string) {
	goalID, err := strconv.Atoi(args[0])
	bail(err)

	err = client.DeleteGoal("", goalID)
	bail(err)
}
