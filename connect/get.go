package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var (
	formatJSON bool
)

func init() {
	getCmd := &cobra.Command{
		Use:   "get <URL>",
		Short: "Get data from Garmin Connect, print to stdout",
		Run:   get,
		Args:  cobra.ExactArgs(1),
	}
	getCmd.Flags().BoolVarP(&formatJSON, "json", "j", false, "Format output as indented JSON")
	rootCmd.AddCommand(getCmd)
}

func get(_ *cobra.Command, args []string) {
	url := args[0]

	if formatJSON {
		raw := bytes.NewBuffer(nil)
		buffer := bytes.NewBuffer(nil)

		err := client.Download(url, raw)
		bail(err)

		err = json.Indent(buffer, raw.Bytes(), "", "  ")
		bail(err)

		_, err = io.Copy(os.Stdout, buffer)
		bail(err)
	} else {
		err := client.Download(url, os.Stdout)
		bail(err)
	}
}
