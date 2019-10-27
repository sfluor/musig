package cmd

import (
	"os"
	"path/filepath"

	"github.com/sfluor/musig/internal/pkg/pipeline"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(loadCmd)
	loadCmd.Flags().BoolP("dry-run", "d", false, "disable saving to the database")
	loadCmd.Flags().BoolP("reset", "r", false, "reset the database if it already exists")
	loadCmd.Flags().BoolP("verbose", "v", false, "enable verbose output")
}

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load [glob]",
	Short: "Load loads all the audio files matching the provided glob into the database (TODO: only .wav are supported for now)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		files, err := filepath.Glob(args[0])
		failIff(err, "error finding files")

		if files == nil {
			log.Infof("no files matched pattern: %s", args[0])
			os.Exit(0)
		}

		resetDB, err := cmd.Flags().GetBool("reset")
		if resetDB && err == nil {
			log.Info("removing the existing database...")
			if err := os.Remove(dbFile); err != nil {
				log.Errorf("Error removing the database at %s: %s", dbFile, err)
			}
		}

		p, err := pipeline.NewDefaultPipeline(dbFile)
		failIff(err, "error creating pipeline")
		defer p.Close()

		process := p.ProcessAndStore

		dryRun, err := cmd.Flags().GetBool("dry-run")
		if dryRun && err == nil {
			log.Info("enabling dry-run (results won't be saved to the database)")
			process = p.Process
		}

		verbose, err := cmd.Flags().GetBool("verbose")
		verbose = verbose && err == nil

		for _, file := range files {
			res, err := process(file)
			failIff(err, "error processing file %s", file)
			if verbose {
				log.Infof("Processed file at %s and got: %+v", file, res.Fingerprint)
			}
		}
	},
}
