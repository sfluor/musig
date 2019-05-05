package db

import (
	"encoding/binary"
	"fmt"

	"github.com/pkg/errors"
	"github.com/sfluor/musig/internal/pkg/model"
	"go.etcd.io/bbolt"
)

var _ Database = &BoltDB{}

// BoltDB implements the Database interface using a bolt database
type BoltDB struct {
	*bbolt.DB
	fingerprintBucket []byte
	songBucket        []byte
}

// NewBoltDB returns a new bolt database
func NewBoltDB(path string) (*BoltDB, error) {
	db, err := bbolt.Open(path, 0666, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating bolt db with path:%s", path)
	}

	boltDB := &BoltDB{
		DB:                db,
		fingerprintBucket: []byte("fingerprint"),
		songBucket:        []byte("song"),
	}

	// Create buckets
	err = db.Update(func(tx *bbolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(boltDB.fingerprintBucket); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists(boltDB.songBucket); err != nil {
			return err
		}
		return nil
	})

	return boltDB, errors.Wrap(err, "error creating buckets")
}

// Get retrieves the provided keys' values from the bolt file
func (db *BoltDB) Get(keys []model.EncodedKey) (map[model.EncodedKey][]model.TableValue, error) {
	res := make(map[model.EncodedKey][]model.TableValue, len(keys))

	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(db.fingerprintBucket)

		for _, k := range keys {
			raw := b.Get(k.Bytes())
			if len(raw) == 0 {
				continue
			}
			val, err := model.ValuesFromBytes(b.Get(k.Bytes()))
			if err != nil {
				return errors.Wrapf(err, "wrong record stored: %v", raw)
			}

			res[k] = val
		}

		return nil
	})

	return res, errors.Wrap(err, "an error occured when reading from bolt")
}

// Set stores the list of (key, value) into the bolt file
func (db *BoltDB) Set(batch map[model.EncodedKey]model.TableValue) error {
	err := db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(db.fingerprintBucket)

		for k, v := range batch {
			// We append the value to the existing array
			// this is because multiple songs can have the same keys
			// TODO: check if the value is not already in the database array to avoid having a same value multiple times
			rawKey := k.Bytes()
			existing := b.Get(rawKey)

			if err := b.Put(rawKey, append(existing, v.Bytes()...)); err != nil {
				return errors.Wrapf(err, "error setting (key: %v, val: %v)", k, v)
			}
		}
		return nil
	})

	return errors.Wrap(err, "an error occured when writing to bolt")
}

// GetSongID does a song name => songID lookup in the database
func (db *BoltDB) GetSongID(name string) (uint32, error) {
	var id uint32

	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(db.songBucket)
		c := b.Cursor()

		// TODO: fixme this is not really efficient
		// Iterate over all the song records
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if string(v) == name {
				id = binary.LittleEndian.Uint32(k)
				return nil
			}
		}

		return fmt.Errorf("could not find id for song name: %s", name)
	})

	return id, errors.Wrap(err, "an error occured when reading from bolt")
}

// GetSong does a songID => song name lookup in the database
func (db *BoltDB) GetSong(songID uint32) (string, error) {
	var name string

	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(db.songBucket)

		raw := b.Get(itob(songID))
		if len(raw) == 0 {
			return fmt.Errorf("got empty song name")
		}

		name = string(raw)

		return nil
	})

	return name, errors.Wrap(err, "an error occured when reading from bolt")
}

// SetSong stores a song name in the database and returns it's song ID
// It first checks if the song is not already stored in the database
func (db *BoltDB) SetSong(song string) (uint32, error) {
	// First check if we haven't the song in the db already
	id, err := db.GetSongID(song)
	if err == nil {
		return id, nil
	}

	var songID uint32
	err = db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(db.songBucket)
		if err != nil {
			return errors.Wrapf(err, "error creating bucket")
		}

		id, err := b.NextSequence()
		if err != nil {
			return errors.Wrap(err, "error getting next sequence")
		}

		songID = uint32(id)
		rawKey := itob(songID)

		return errors.Wrap(b.Put(rawKey, []byte(song)), "error setting song")
	})

	return songID, errors.Wrap(err, "an error occured when writing to bolt")
}

func itob(s uint32) []byte {
	raw := make([]byte, 4)
	binary.LittleEndian.PutUint32(raw, s)
	return raw
}
