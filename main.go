package main

import (
	"flag"
	"fmt"
	"image/png"
	"os"

	"github.com/sfluor/musig/dsp"
	"github.com/sfluor/musig/fingerprint"
	"github.com/sfluor/musig/model"
)

func main() {
	audioFile := flag.String("file", "", "Audio file to process")
	printFingerprint := flag.Bool("fingerprint", true, "enable / disable printing the fingerprint")
	specFile := flag.String("spec_file", "", "File where we should save the spectrogram (if not specified the spectrogram won't be saved)")

	flag.Parse()

	file, err := os.Open(*audioFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	s := dsp.NewSpectrogrammer(
		model.DOWNSAMPLERATIO,
		model.MAXFREQ,
		model.SAMPLESIZE,
	)

	spec, spr, err := s.Spectrogram(file)
	if err != nil {
		panic(err)
	}

	if *printFingerprint {
		cMap := s.ConstellationMap(spec, spr)
		fpr := fingerprint.NewDefaultFingerprinter()
		// TODO change ID
		songFpr := fpr.Fingerprint(1, cMap)
		fmt.Println("Fingerprint for the given song:")
		for key, val := range songFpr {
			fmt.Printf("key: %v, val: %v\n", key, val)
		}
	}

	if *specFile != "" {
		img := dsp.SpecToImg(spec)
		fmt.Printf("Saving spectrogram to file: %s\n", *specFile)

		f, err := os.Create(*specFile)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if err := png.Encode(f, img); err != nil {
			panic(err)
		}
	}
}
