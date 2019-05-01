package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(readCmd)
}

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read [file]",
	Short: "Read reads the given audio file trying to find it's song name",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("read called")
	},
}
