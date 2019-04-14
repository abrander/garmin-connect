package main

import (
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"

	connect "github.com/abrander/garmin-connect"
)

var (
	rootCmd = &cobra.Command{
		Use:   os.Args[0] + " [command]",
		Short: "CLI Client for Garmin Connect",
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
		Use:   "authenticate [email]",
		Short: "Authenticate against the Garmin API",
		Run:   authenticate,
		Args:  cobra.RangeArgs(0, 1),
	}
	rootCmd.AddCommand(authenticateCmd)

	signoutCmd := &cobra.Command{
		Use:   "signout",
		Short: "Log out of the Garmin API and forget session and password",
		Run:   signout,
		Args:  cobra.NoArgs,
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

func authenticate(_ *cobra.Command, args []string) {
	var email string
	if len(args) == 1 {
		email = args[0]
	} else {
		fmt.Print("Email: ")
		fmt.Scanln(&email)
	}

	fmt.Print("Password: ")

	password, err := terminal.ReadPassword(syscall.Stdin)
	bail(err)

	client.SetOptions(connect.Credentials(email, string(password)))
	err = client.Authenticate()
	bail(err)

	fmt.Printf("\nSuccess\n")
}

func signout(_ *cobra.Command, _ []string) {
	_ = client.Signout()
	client.Password = ""
}
