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

		for _, sample := range samples {
			fmt.Printf("%f,%d\n", time, sample)
			time += 1 / float64(reader.SampleRate())
		}
	}
}
