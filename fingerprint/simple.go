package fingerprint

import (
	"github.com/sfluor/musig/model"
)

var _ Fingerprinter = &SimpleFingerprinter{}

// SimpleFingerprinter is a fingerprinter that just takes
type SimpleFingerprinter struct {
	anchorOffset   int
	targetZoneSize int
	lastOffset     int
}

// NewDefaultFingerprinter returns a new fingerprinter using the simple strategy
// using an anchor offset of 2 and target zone size of 5
func NewDefaultFingerprinter() *SimpleFingerprinter {
	return NewSimpleFingerprinter(2, 5)
}

// NewSimpleFingerprinter returns a new simple fingerprinter
func NewSimpleFingerprinter(anchorOffset, targetZoneSize int) *SimpleFingerprinter {
	return &SimpleFingerprinter{
		anchorOffset:   anchorOffset,
		targetZoneSize: targetZoneSize,
		lastOffset:     anchorOffset + targetZoneSize,
	}
}

func (sf *SimpleFingerprinter) Fingerprint(songID uint32, cMap []model.ConstellationPoint) map[model.EncodedKey]model.TableValue {
	length := len(cMap)
	res := map[model.EncodedKey]model.TableValue{}

	for i := 0; i+sf.lastOffset < length; i++ {
		anchor := cMap[i]
		for _, p := range cMap[i+sf.anchorOffset : i+sf.lastOffset] {
			res[model.NewTableKey(anchor, p).Encode()] = *model.NewTableValue(songID, anchor)
		}
	}

	return res
}
