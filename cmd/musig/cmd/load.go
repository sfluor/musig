package cmd

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sfluor/musig/internal/pkg/db"
	"github.com/sfluor/musig/internal/pkg/fingerprint"
	"github.com/sfluor/musig/internal/pkg/model"
	"github.com/sfluor/musig/pkg/dsp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(loadCmd)
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

		db, err := db.NewBoltDB(dbFile)
		failIff(err, "error connection to database at: %s", dbFile)
		defer db.Close()

		for _, file := range files {
			err := loadFile(db, file)
			failIff(err, "error loading file %s", file)
		}
	},
}

// loadFile loads the given file in the database
func loadFile(db db.Database, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return errors.Wrapf(err, "error opening file at %s", path)
	}
	defer file.Close()

	s := dsp.NewSpectrogrammer(model.DOWNSAMPLERATIO, model.MAXFREQ, model.SAMPLESIZE)

	spec, spr, err := s.Spectrogram(file)
	if err != nil {
		return errors.Wrap(err, "error generating spectrogram")
	}

	cMap := s.ConstellationMap(spec, spr)
	fpr := fingerprint.NewDefaultFingerprinter()

	id, err := db.SetSong(path)
	if err != nil {
		return errors.Wrap(err, "error storing song name in database")
	}

	songFpr := fpr.Fingerprint(id, cMap)
	if err := db.Set(songFpr); err != nil {
		return errors.Wrap(err, "error storing song fingerprint in database")
	}

	log.Infof("sucessfully loaded %s into the database", path)
	return nil
}
