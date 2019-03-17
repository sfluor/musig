package main

import (
	"fmt"
	"image/png"
	"io"
	"os"

	"github.com/sfluor/musig/dsp"
	"github.com/sfluor/musig/sound"
)

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader, err := sound.NewWAVReader(file)
	if err != nil {
		panic(err)
	}

	binSize := float64(1024)

	lp := dsp.NewLPFilter(dsp.MAXFREQ, reader.SampleRate())

	// TIME : FREQUENCY : Value
	// time is t * dsp.DOWNSAMPLERATIO / reader.SampleRate()
	// frequency is f * freqBinSize
	matrix := [][]float64{}

	bin := make([]float64, int(binSize*dsp.DOWNSAMPLERATIO))
	for {
		n, err := reader.Read(bin, int(binSize*dsp.DOWNSAMPLERATIO))
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintf(os.Stderr, "error reading from sound file: %s\n", err)
			os.Exit(1)
		}

		fft := dsp.FFT(
			dsp.Downsample(
				lp.Filter(bin[:n]),
				dsp.DOWNSAMPLERATIO,
			),
		)

		// TODO remove slicing here when removing duplicate values retuned by the FFT
		matrix = append(matrix, fft[:len(fft)-(len(fft)-1)/2])
	}

	spec := dsp.Spectrogram(matrix)
	f, err := os.Create(os.Args[2])
	if err != nil {
		panic(err)
	}

	if err := png.Encode(f, spec); err != nil {
		panic(err)
	}
}
