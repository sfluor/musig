package pipeline

import (
	"os"

	"github.com/pkg/errors"
	"github.com/sfluor/musig/internal/pkg/db"
	"github.com/sfluor/musig/internal/pkg/fingerprint"
	"github.com/sfluor/musig/internal/pkg/model"
	"github.com/sfluor/musig/pkg/dsp"
	log "github.com/sirupsen/logrus"
)

// Pipeline is a struct that allows to operate on audio files
type Pipeline struct {
	s   *dsp.Spectrogrammer
	DB  db.Database
	fpr fingerprint.Fingerprinter
}

// NewDefaultPipeline creates a new default pipeline using a Bolt DB, the default fingerprinter and the default spectrogrammer
func NewDefaultPipeline(dbFile string) (*Pipeline, error) {
	db, err := db.NewBoltDB(dbFile)
	if err != nil {
		return nil, errors.Wrapf(err, "error connection to database at: %s", dbFile)
	}
	s := dsp.NewSpectrogrammer(model.DownsampleRatio, model.MaxFreq, model.SampleSize)
	fpr := fingerprint.NewDefaultFingerprinter()

	return &Pipeline{
		s:   s,
		DB:  db,
		fpr: fpr,
	}, nil
}

// NewPipeline creates a new pipeline
func NewPipeline(s *dsp.Spectrogrammer, db db.Database, fpr fingerprint.Fingerprinter) *Pipeline {
	return &Pipeline{
		s:   s,
		DB:  db,
		fpr: fpr,
	}
}

// Close closes the underlying database
func (p *Pipeline) Close() {
	p.DB.Close()
}

// Result represents the output of a pipeline
type Result struct {
	Path         string
	CMap         []model.ConstellationPoint
	SongID       uint32
	SamplingRate float64
	Spectrogram  [][]float64
	Fingerprint  map[model.EncodedKey]model.TableValue
}

// ProcessAndStore process the given audio file and store it in the database
// the computed results are returned
func (p *Pipeline) ProcessAndStore(path string) (*Result, error) {
	partial, err := p.read(path)
	if err != nil {
		return nil, err
	}

	id, err := p.DB.SetSong(model.SongNameFromPath(path))
	if err != nil {
		return nil, errors.Wrap(err, "error storing song name in database")
	}

	songFpr := p.fpr.Fingerprint(id, partial.cMap)
	if err := p.DB.Set(songFpr); err != nil {
		return nil, errors.Wrap(err, "error storing song fingerprint in database")
	}

	log.Infof("sucessfully loaded %s into the database", path)
	return &Result{
		Path:         path,
		CMap:         partial.cMap,
		SongID:       id,
		SamplingRate: partial.samplingRate,
		Spectrogram:  partial.spectrogram,
		Fingerprint:  songFpr,
	}, nil
}

// Process process the given audio file and returns a Result
func (p *Pipeline) Process(path string) (*Result, error) {
	partial, err := p.read(path)
	if err != nil {
		return nil, err
	}

	// Use 0 as the ID
	var id uint32
	songFpr := p.fpr.Fingerprint(id, partial.cMap)

	log.Infof("sucessfully loaded %s into the database", path)
	return &Result{
		Path:         path,
		CMap:         partial.cMap,
		SongID:       id,
		SamplingRate: partial.samplingRate,
		Spectrogram:  partial.spectrogram,
		Fingerprint:  songFpr,
	}, nil
}

type partialResult struct {
	cMap         []model.ConstellationPoint
	samplingRate float64
	spectrogram  [][]float64
}

func (p *Pipeline) read(path string) (*partialResult, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "error opening file at %s", path)
	}
	defer file.Close()

	spec, spr, err := p.s.Spectrogram(file)
	if err != nil {
		return nil, errors.Wrap(err, "error generating spectrogram")
	}

	cMap := p.s.ConstellationMap(spec, spr)

	return &partialResult{
		cMap:         cMap,
		samplingRate: spr,
		spectrogram:  spec,
	}, nil
}
