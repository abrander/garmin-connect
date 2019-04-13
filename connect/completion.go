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
		Use:   "bash",
		Short: "Output command completion for Bourne Again Shell (bash)",
		RunE:  completionBash,
		Args:  cobra.NoArgs,
	}
	completionCmd.AddCommand(completionBashCmd)

	completionZshCmd := &cobra.Command{
		Use:   "zsh",
		Short: "Output command completion for Z Shell (zsh)",
		RunE:  completionZsh,
		Args:  cobra.NoArgs,
	}
	completionCmd.AddCommand(completionZshCmd)
}

func completionBash(_ *cobra.Command, _ []string) error {
	return rootCmd.GenBashCompletion(os.Stdout)
}

func completionZsh(_ *cobra.Command, _ []string) error {
	return rootCmd.GenZshCompletion(os.Stdout)
}
