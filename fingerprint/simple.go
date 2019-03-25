package fingerprint

import (
	"fmt"

	"github.com/sfluor/musig/model"
)

var _ Fingerprinter = &SimpleFingerprinter{}

// SimpleFingerprinter is a fingerprinter that just takes
type SimpleFingerprinter struct {
	anchorOffset   int
	targetZoneSize int
	lastOffset     int
}

// NewSimpleFingerprinter returns a new simple fingerprinter
func NewSimpleFingerprinter(anchorOffset, targetZoneSize int) *SimpleFingerprinter {
	return &SimpleFingerprinter{
		anchorOffset:   anchorOffset,
		targetZoneSize: targetZoneSize,
		lastOffset:     anchorOffset + targetZoneSize,
	}
}

func (sf *SimpleFingerprinter) Fingerprint(songID uint32, cMap []model.ConstellationPoint) map[model.TableKey]model.TableValue {
	length := len(cMap)
	res := map[model.TableKey]model.TableValue{}

	for i := 0; i+sf.lastOffset < length; i++ {
		anchor := cMap[i]
		for _, p := range cMap[i+sf.anchorOffset : i+sf.lastOffset] {
			res[*model.NewTableKey(anchor, p)] = *model.NewTableValue(songID, anchor)
		}
	}

	fmt.Printf("res = %+v\n", res)
	return res
}
