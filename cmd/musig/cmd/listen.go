package cmd

import (
	"os"
	"path"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listenCmd)
	listenCmd.Flags().DurationP("duration", "d", 10*time.Second, "duration of the listening")
}

// listenCmd represents the listen command
var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "listen will record the microphone input and try to find a matching song from the database (Ctrl-C will stop the recording)",
	Run: func(cmd *cobra.Command, args []string) {
		name := path.Join(os.TempDir(), "musig_record.wav")
		dur, err := cmd.Flags().GetDuration("duration")
		failIff(err, "could not get duration, got: %v", dur)

		recordAudioToFile(name, dur)
		cmdRead(name)
		// todo remove
	},
}
