package main

import (
	"image/png"
	"os"

	"github.com/sfluor/musig/dsp"
)

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	s := dsp.NewSpectrogrammer(
		dsp.DOWNSAMPLERATIO,
		dsp.MAXFREQ,
		1024,
	)

	spec, _, err := s.Spectrogram(file)
	if err != nil {
		panic(err)
	}

	img := dsp.SpecToImg(spec)

	f, err := os.Create(os.Args[2])
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}
