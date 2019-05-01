package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var dbFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "musig",
	Short: "A shazam like CLI tool",
	Long:  `TODO`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&dbFile, "database", "/tmp/bolt.db", "database file to use")
}

func failIff(err error, msg string, args ...interface{}) {
	if err != nil {
		fmt.Fprintf(os.Stderr, msg, args...)
		os.Exit(1)
	}
}
