package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path"

	"github.com/mitchellh/go-homedir"

	"github.com/abrander/garmin-connect"
)

var state struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	SessionID string `json:"sessionID"`
}

var client *connect.Client

func stateFilename() string {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatalf("Could not detect home directory: %s", err.Error())
	}

	return path.Join(home, ".garmin-connect.json")
}

func loadState() {
	client = connect.NewClient(
		connect.AutoRenewSession(true),
	)
	data, err := ioutil.ReadFile(stateFilename())
	if err != nil {
		log.Printf("Could not open state file: %s", err.Error())
		return
	}

	err = json.Unmarshal(data, &state)
	if err != nil {
		log.Fatalf("Could not unmarshal state: %s", err.Error())
	}

	client.SetOptions(
		connect.Credentials(state.Email, state.Password),
		connect.SessionID(state.SessionID),
	)
}

func storeState() {
	if client != nil {
		state.SessionID = client.SessionID()
	}

	b, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		log.Fatalf("Could not marshal state: %s", err.Error())
	}

	err = ioutil.WriteFile(stateFilename(), b, 0600)
	if err != nil {
		log.Fatalf("Could not write state file: %s", err.Error())
	}
}
