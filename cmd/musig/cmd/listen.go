package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listenCmd)
}

// listenCmd represents the listen command
var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "TODO",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO")
	},
}
