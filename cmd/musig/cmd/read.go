package cmd

import (
	"fmt"
	"github.com/sfluor/musig/internal/pkg/model"
	"github.com/sfluor/musig/internal/pkg/pipeline"
	"github.com/sfluor/musig/pkg/stats"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(readCmd)
}

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read [file]",
	Short: "Read reads the given audio file trying to find it's song name",
	Run: func(cmd *cobra.Command, args []string) {
		cmdRead(args[0])
	},
}

// TODO move this outside
func correlation(sample, match map[model.EncodedKey]model.TableValue) float64 {
	// Will hold a list of points (time in the sample sound file, time in the matched database sound file)
	points := [2][]float64{}
	for k, sampleValue := range sample {
		if matchValue, ok := match[k]; ok {
			points[0] = append(points[0], float64(sampleValue.AnchorTimeMs))
			points[1] = append(points[1], float64(matchValue.AnchorTimeMs))
		}
	}
	return stats.Correlation(points[0], points[1])
}

func cmdRead(file string) {
	p, err := pipeline.NewDefaultPipeline(dbFile)
	failIff(err, "error creating pipeline")
	defer p.Close()

	res, err := p.Process(file)
	failIff(err, "error processing file %s", file)

	// Will hold a count of songID => occurrences
	counts := map[uint32]int{}
	keys := make([]model.EncodedKey, 0, len(res.Fingerprint))
	sample := map[model.EncodedKey]model.TableValue{}
	// songID => points that matched
	matches := map[uint32]map[model.EncodedKey]model.TableValue{}

	for k, v := range res.Fingerprint {
		keys = append(keys, k)
		sample[k] = v
	}

	m, err := p.DB.Get(keys)
	for key, values := range m {
		for _, val := range values {
			counts[val.SongID] += 1

			if _, ok := matches[val.SongID]; !ok {
				matches[val.SongID] = map[model.EncodedKey]model.TableValue{}
			}

			matches[val.SongID][key] = val
		}
	}

	// songID => correlation
	correlations := map[uint32]float64{}
	for songID, points := range matches {
		correlations[songID] = correlation(sample, points)
	}

	var song string
	var max, total int
	fmt.Println("Matches:")
	for id, count := range counts {
		name, err := p.DB.GetSong(id)
		failIff(err, "error getting song id: %d", id)
		fmt.Printf("\t- %s, count: %d, correlation: %f\n", name, count, correlations[id])
		if count > max {
			song, max = name, count
		}
		total += count
	}

	fmt.Println("---")
	fmt.Printf("Song is: %s (count: %d, pct: %.2f %%)\n", song, max, 100*float64(max)/float64(total))
}
