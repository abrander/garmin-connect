package main

import (
	"os"

	"github.com/spf13/cobra"
)

func init() {
	completionCmd := &cobra.Command{
		Use: "completion",
	}
	rootCmd.AddCommand(completionCmd)

	completionBashCmd := &cobra.Command{
		Use: "bash",
		Run: completionBash,
	}
	completionCmd.AddCommand(completionBashCmd)

	completionZshCmd := &cobra.Command{
		Use: "zsh",
		Run: completionZsh,
	}
	completionCmd.AddCommand(completionZshCmd)
}

func completionBash(_ *cobra.Command, _ []string) {
	rootCmd.GenBashCompletion(os.Stdout)
}

func completionZsh(_ *cobra.Command, _ []string) {
	rootCmd.GenZshCompletion(os.Stdout)
}
