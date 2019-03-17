package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"os"

	"github.com/pkg/errors"
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

	if err := matrixToImg(os.Args[2], matrix); err != nil {
		panic(err)
	}
}

func matrixToImg(file string, matrix [][]float64) error {
	// 10 Pixel for 1 frequency step
	// 150 Pixels for 1 timeStep
	nTime, nFreq := 150, 10
	img := image.NewRGBA(image.Rect(0, 0, nTime*len(matrix), nFreq*len(matrix[0])))

	// TODO change those hard coded values
	min, max := -1000.0, 1000.0
	for t, row := range matrix {
		for f, a := range row {
			c := colorbar((a - min) / (max - min))
			fmt.Printf("c = %+v\n", c)
			for i := 0; i < nTime; i++ {
				for j := 0; j < nFreq; j++ {
					img.Set(nTime*t+i, nFreq*len(matrix[0])-nFreq*f+j, c)
				}
			}
		}
	}

	f, err := os.Create(file)
	if err != nil {
		return errors.Wrapf(err, "error creating file %s", file)
	}
	return png.Encode(f, img)
}

func colorbar(val float64) color.RGBA {
	if val == 0 {
		return color.RGBA{}
	}

	r := 255 * math.Min(math.Max(0, 1.5-math.Abs(1-4*(val-0.5))), 1)
	g := 255 * math.Min(math.Max(0, 1.5-math.Abs(1-4*(val-0.25))), 1)
	b := 255 * math.Min(math.Max(0, 1.5-math.Abs(1-4*val)), 1)
	return color.RGBA{uint8(r), uint8(g), uint8(b), 255}
}
