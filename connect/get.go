package main

import (
	"os"

	"github.com/spf13/cobra"
)

func init() {
	getCmd := &cobra.Command{
		Use:   "get <URL>",
		Short: "Get data from Garmin Connect, print to stdout",
		Run:   get,
		Args:  cobra.ExactArgs(1),
	}
	rootCmd.AddCommand(getCmd)
}

func get(_ *cobra.Command, args []string) {
	url := args[0]

	err := client.Download(url, os.Stdout)
	bail(err)
}
