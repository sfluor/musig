package cmd

import (
	"io/ioutil"
	"os"
	"os/signal"
	"time"

	"github.com/sfluor/musig/internal/pkg/sound"
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
		tmpFile, err := ioutil.TempFile(os.TempDir(), "musig_record.wav")
		failIff(err, "error creating temporary file for recording in %s", tmpFile.Name())
		defer os.Remove(tmpFile.Name())

		stopCh := make(chan struct{}, 1)
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, os.Kill)

		dur, err := cmd.Flags().GetDuration("duration")
		failIff(err, "could not get duration, got: %v", dur)

		go func() {
			defer func() { stopCh <- struct{}{} }()
			for {
				select {
				case <-time.After(dur):
					return
				case <-sig:
					return
				}
			}
		}()

		err = sound.RecordWAV(tmpFile, stopCh)
		failIff(err, "an error occured recording WAV file")
		failIff(tmpFile.Sync(), "error syncing temp file")

		cmdRead(tmpFile.Name())
	},
}
