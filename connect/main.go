package main

import (
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/abrander/garmin-connect"
)

var (
	rootCmd = &cobra.Command{
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			loadState()
			if verbose {
				logger := log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
				client.SetOptions(connect.DebugLogger(logger))
			}

			if dumpFile != "" {
				w, err := os.OpenFile(dumpFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
				bail(err)
				client.SetOptions(connect.DumpWriter(w))
			}
		},
		PersistentPostRun: func(_ *cobra.Command, _ []string) {
			storeState()
		},
	}

	verbose  bool
	dumpFile string
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose debug output")
	rootCmd.PersistentFlags().StringVarP(&dumpFile, "dump", "d", "", "File to dump requests and responses to")

	authenticateCmd := &cobra.Command{
		Use: "authenticate",
		Run: authenticate,
	}
	rootCmd.AddCommand(authenticateCmd)

	signoutCmd := &cobra.Command{
		Use: "signout",
		Run: signout,
	}
	rootCmd.AddCommand(signoutCmd)
}

func bail(err error) {
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
}

func main() {
	bail(rootCmd.Execute())
}

func authenticate(_ *cobra.Command, _ []string) {
	var email string
	fmt.Print("Email: ")
	fmt.Scanln(&email)
	fmt.Print("Password: ")

	password, err := terminal.ReadPassword(int(syscall.Stdin))
	bail(err)

	client.SetOptions(connect.Credentials(string(email), string(password)))
	err = client.Authenticate()
	bail(err)

	state.Email = string(email)
	state.Password = string(password)
	state.SessionID = client.SessionID()

	fmt.Printf("\nSuccess\n")
}

func signout(_ *cobra.Command, _ []string) {
	client.Signout()

	state.Email = ""
	state.Password = ""
	state.SessionID = ""
	client = nil

	storeState()
}
