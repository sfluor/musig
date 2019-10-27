package cmd

import (
	"fmt"
	"github.com/sfluor/musig/internal/pkg/model"
	"github.com/sfluor/musig/internal/pkg/pipeline"
	"github.com/sfluor/musig/pkg/dsp"
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

func cmdRead(file string) {
	p, err := pipeline.NewDefaultPipeline(dbFile)
	failIff(err, "error creating pipeline")
	defer p.Close()

	res, err := p.Process(file)
	failIff(err, "error processing file %s", file)

	// Will hold a count of songID => occurrences
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

			if _, ok := matches[val.SongID]; !ok {
				matches[val.SongID] = map[model.EncodedKey]model.TableValue{}
			}

			matches[val.SongID][key] = val
		}
	}

	// songID => correlation
	scores := map[uint32]float64{}
	for songID, points := range matches {
		scores[songID] = dsp.MatchScore(sample, points)
	}

	var song string
	var max float64
	fmt.Println("Matches:")
	for id, score := range scores {
		name, err := p.DB.GetSong(id)
		failIff(err, "error getting song id: %d", id)
		fmt.Printf("\t- %s, score: %f\n", name, score)
		if score > max {
			song, max = name, score
		}
	}

	fmt.Println("---")
	fmt.Printf("Song is: %s (score: %f)\n", song, max)
}
