package model

import (
	"fmt"
	"math"
)

// Used to down size the frequences to a 9 bit ints
// XXX: we could also use 10.7Hz directly
const freqStep = MAXFREQ / float64(1<<9)

// Used to down size the delta times to 14 bit ints (we use 16s as the max duration)
const deltaTimeStep = 16 / float64(1<<14)

// TableKey represents a table key
type TableKey struct {
	// Frequency of the anchor point for the given point's target zone
	AnchorFreq float64
	// Frequency of the given point
	PointFreq float64
	// Delta time between the anchor point and the given point
	DeltaT float64
}

// EncodedKey represents an encoded key
type EncodedKey uint32

func NewTableKey(anchor, point ConstellationPoint) *TableKey {
	return &TableKey{
		AnchorFreq: anchor.Freq,
		PointFreq:  point.Freq,
		// Use absolute just in case anchor is after the target zone
		DeltaT: math.Abs(point.Time - anchor.Time),
	}
}

// Encode encodes the table key using:
// 9 bits for the “frequency of the anchor”: fa
// 9 bits for the ” frequency of the point”: fp
// 14 bits for the ”delta time between the anchor and the point”: dt
// The result is then dt | fa | fp
// XXX: this only works if frequencies are coded in 9 bits or less (if we used a 1024 samples FFT, it will be the case)
func (tk *TableKey) Encode() EncodedKey {
	// down size params
	fp := uint32(tk.PointFreq / freqStep)
	fa := uint32(tk.AnchorFreq / freqStep)
	dt := uint32(tk.DeltaT / deltaTimeStep)

	res := fp
	res |= fa << 9
	res |= dt << 23
	return EncodedKey(res)
}

// TableValue represents a table value
type TableValue struct {
	// AnchorTimeMs is the time of the anchor in the related song in milliseconds
	AnchorTimeMs uint32
	// SongID is an ID representing the related song
	SongID uint32
}

func NewTableValue(song uint32, anchor ConstellationPoint) *TableValue {
	return &TableValue{
		AnchorTimeMs: uint32(anchor.Time * 1000),
		SongID:       song,
	}
}

func (tv TableValue) String() string {
	return fmt.Sprintf("(anchor_time_ms: %d, song_id: %d)", tv.AnchorTimeMs, tv.SongID)
}
