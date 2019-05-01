package cmd

import (
	"image/png"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"
	"github.com/sfluor/musig/internal/pkg/model"
	"github.com/sfluor/musig/pkg/dsp"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(specCmd)
}

// specCmd represents the spectrogram command
var specCmd = &cobra.Command{
	Use:   "spectrogram [audio_file] [output_img]",
	Short: "spectrogram generate a spectrogram image for the given audio file in png (TODO: only .wav are supported for now)",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		err := genSpectrogram(args[0], args[1])
		failIff(err, "error generating spectrogram")
	},
}

// genSpectrogram generates a spectrogram for the given file
func genSpectrogram(path string, imgPath string) error {
	log.Infof("reading %s to save it's spectrogram to: %s", path, imgPath)
	file, err := os.Open(path)
	if err != nil {
		return errors.Wrapf(err, "error opening file at %s", path)
	}
	defer file.Close()

	s := dsp.NewSpectrogrammer(model.DOWNSAMPLERATIO, model.MAXFREQ, model.SAMPLESIZE)

	spec, _, err := s.Spectrogram(file)
	if err != nil {
		return errors.Wrap(err, "error generating spectrogram")
	}

	img := dsp.SpecToImg(spec)

	f, err := os.Create(imgPath)
	if err != nil {
		return errors.Wrapf(err, "error creating image file at %s", imgPath)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		return errors.Wrap(err, "error encoding image file")
	}
	log.Infof("successfully saved spectrogram to file: %s", imgPath)

	return nil
}
