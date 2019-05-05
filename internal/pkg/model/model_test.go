package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSongNameFromPath(t *testing.T) {
	for _, tc := range []struct {
		path     string
		expected string
	}{
		{
			path:     "assets/dataset/wav/Summer Was Fun & Laura Brehm - Prism [NCS Release].wav",
			expected: "Summer Was Fun & Laura Brehm - Prism [NCS Release]",
		},
		{
			path:     " assets/dataset/wav/BEAUZ & Momo - Won't Look Back [NCS Release].wav ",
			expected: "BEAUZ & Momo - Won't Look Back [NCS Release]",
		},
		{
			path:     "assets/dataset/wav/Maduk - Go Home (Original Mix) [NCS Release].wav",
			expected: "Maduk - Go Home (Original Mix) [NCS Release]",
		},
		{
			path:     "./test.mp3",
			expected: "test",
		},
	} {
		assert.Equal(t, tc.expected, SongNameFromPath(tc.path))
	}
}
