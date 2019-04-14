package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"

	connect "github.com/abrander/garmin-connect"
)

var (
	client = connect.NewClient(
		connect.AutoRenewSession(true),
	)

	stateFile string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&stateFile, "state", "s", stateFilename(), "State file to use")
}

func stateFilename() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Could not detect home directory: %s", err.Error())
	}

	return path.Join(home, ".garmin-connect.json")
}

func loadState() {
	data, err := ioutil.ReadFile(stateFile)
	if err != nil {
		log.Printf("Could not open state file: %s", err.Error())
		return
	}

	err = json.Unmarshal(data, client)
	if err != nil {
		log.Fatalf("Could not unmarshal state: %s", err.Error())
	}
}

func storeState() {
	b, err := json.MarshalIndent(client, "", "  ")
	if err != nil {
		log.Fatalf("Could not marshal state: %s", err.Error())
	}

	err = ioutil.WriteFile(stateFile, b, 0600)
	if err != nil {
		log.Fatalf("Could not write state file: %s", err.Error())
	}
}
