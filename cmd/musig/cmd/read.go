package cmd

import (
	"fmt"

	"github.com/sfluor/musig/internal/pkg/model"
	"github.com/sfluor/musig/internal/pkg/pipeline"
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
		p, err := pipeline.NewDefaultPipeline(dbFile)
		failIff(err, "error creating pipeline")
		defer p.Close()

		res, err := p.Process(args[0])
		failIff(err, "error processing file %s", args[0])

		// Will hold a count of songID => occurences
		counts := map[uint32]int{}
		keys := make([]model.EncodedKey, 0, len(res.Fingerprint))

		for k := range res.Fingerprint {
			keys = append(keys, k)
		}

		m, err := p.DB.Get(keys)
		for _, v := range m {
			counts[v.SongID] += 1
		}

		var song string
		var max, total int
		fmt.Println("Matches:")
		for id, count := range counts {
			name, err := p.DB.GetSong(id)
			failIff(err, "error getting song id: %d", id)
			fmt.Printf("\t- %s, count: %d\n", name, count)
			if count > max {
				song, max = name, count
			}
			total += count
		}

		fmt.Println("---")
		fmt.Printf("Song is: %s (count: %d, pct: %.2f %%)\n", song, max, 100*float64(max)/float64(total))
	},
}
