package main

import (
	"fmt"
	"io"
	"os"

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

	var time float64
	for {
		samples, err := reader.Read()
		if err == io.EOF {
			break
		}

		lp := sound.NewLPFilter(sound.MAXFREQ, reader.SampleRate())
		filtered := sound.Downsample(lp.Filter(samples), sound.DOWNSAMPLERATIO)

		for _, sample := range filtered {
			fmt.Printf("%f,%f\n", time, sample)
			time += 1 / reader.SampleRate()
		}
	}
}
