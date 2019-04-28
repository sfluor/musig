package fingerprint

import "github.com/sfluor/musig/internal/pkg/model"

type Fingerprinter interface {
	Fingerprint(uint32, []model.ConstellationPoint) map[model.EncodedKey]model.TableValue
}
