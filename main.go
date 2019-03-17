package main

import (
	"fmt"
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
	// freqBinSize := reader.SampleRate() / dsp.DOWNSAMPLERATIO / binSize

	lp := dsp.NewLPFilter(dsp.MAXFREQ, reader.SampleRate())

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
		fmt.Printf("fft = %+v\n", fft)
		fmt.Printf("len(fft) = %+v\n", len(fft))
	}
}
