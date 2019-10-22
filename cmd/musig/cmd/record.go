package cmd

import (
	"os"
	"os/signal"
	"time"

	"github.com/sfluor/musig/internal/pkg/sound"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(recordCmd)
	recordCmd.Flags().DurationP("duration", "d", 10*time.Second, "duration of the listening")
}

// recordCmd represents the listen command
var recordCmd = &cobra.Command{
	Use:   "record",
	Short: "record will record the microphone input and save the signal to the given file",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dur, err := cmd.Flags().GetDuration("duration")
		failIff(err, "could not get duration, got: %v", dur)

		recordAudioToFile(args[0], dur)
	},
}

func recordAudioToFile(name string, duration time.Duration) {
	file, err := os.Create(name)
	failIff(err, "error creating file for recording in %s", name)

	stopCh := make(chan struct{}, 1)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	go func() {
		defer func() {
			stopCh <- struct{}{}
		}()
		for {
			select {
			case <-time.After(duration):
				return
			case <-sig:
				return
			}
		}
	}()

	err = sound.RecordWAV(file, stopCh)
	failIff(err, "an error occurred recording WAV file")
	failIff(file.Sync(), "error syncing temp file")
}
