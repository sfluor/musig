package model

import (
	"encoding/binary"
	"fmt"
	"math"
)

// Used to down size the frequences to a 9 bit ints
// XXX: we could also use 10.7Hz directly
const freqStep = MaxFreq / float64(1<<9)

// Used to down size the delta times to 14 bit ints (we use 16s as the max duration)
const deltaTimeStep = 16 / float64(1<<14)

// TableValueSize represents the TableValueSize when encoded in bytes
const TableValueSize = 8

// AnchorKey represents a anchor key
type AnchorKey struct {
	// Frequency of the anchor point for the given point's target zone
	AnchorFreq float64
	// Frequency of the given point
	PointFreq float64
	// Delta time between the anchor point and the given point
	DeltaT float64
}

// EncodedKey represents an encoded key
type EncodedKey uint32

// NewAnchorKey creates a new anchor key from the given anchor and the given point
func NewAnchorKey(anchor, point ConstellationPoint) *AnchorKey {
	return &AnchorKey{
		AnchorFreq: anchor.Freq,
		PointFreq:  point.Freq,
		// Use absolute just in case anchor is after the target zone
		DeltaT: math.Abs(point.Time - anchor.Time),
	}
}

// Bytes encodes the key in bytes
func (ek EncodedKey) Bytes() []byte {
	// uint32 is 4 bytes
	bk := make([]byte, 4)
	binary.LittleEndian.PutUint32(bk, uint32(ek))
	return bk
}

// Encode encodes the anchor key using:
// 9 bits for the “frequency of the anchor”: fa
// 9 bits for the ” frequency of the point”: fp
// 14 bits for the ”delta time between the anchor and the point”: dt
// The result is then dt | fa | fp
// XXX: this only works if frequencies are coded in 9 bits or less (if we used a 1024 samples FFT, it will be the case)
func (tk *AnchorKey) Encode() EncodedKey {
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

// NewTableValue creates a new table value from the given song ID and anchor point
func NewTableValue(song uint32, anchor ConstellationPoint) *TableValue {
	return &TableValue{
		AnchorTimeMs: uint32(anchor.Time * 1000),
		SongID:       song,
	}
}

// Bytes encodes the given table value in bytes
func (tv *TableValue) Bytes() []byte {
	// Use a uint64 (8 bytes)
	b := make([]byte, TableValueSize)
	binary.LittleEndian.PutUint32(b[:4], tv.AnchorTimeMs)
	binary.LittleEndian.PutUint32(b[4:], tv.SongID)
	return b
}

// ValuesFromBytes decodes a list of table values from the given byte array
func ValuesFromBytes(b []byte) ([]TableValue, error) {
	if len(b)%TableValueSize != 0 {
		return nil, fmt.Errorf("error wrong size for value: %d (got: %v) expected a multiple of %d", len(b), b, TableValueSize)
	}

	N := len(b) / TableValueSize
	res := make([]TableValue, 0, N)

	for i := 0; i < N; i++ {
		tv := TableValue{}
		tv.AnchorTimeMs = binary.LittleEndian.Uint32(b[i*8 : i*8+4])
		tv.SongID = binary.LittleEndian.Uint32(b[i*8+4 : (i+1)*8])
		res = append(res, tv)
	}

	return res, nil
}

// String returns a string representation of a TableValue
func (tv TableValue) String() string {
	return fmt.Sprintf("(anchor_time_ms: %d, song_id: %d)", tv.AnchorTimeMs, tv.SongID)
}
