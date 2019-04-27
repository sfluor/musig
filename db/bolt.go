package db

import (
	"github.com/pkg/errors"
	"github.com/sfluor/musig/model"
	"go.etcd.io/bbolt"
)

var _ Database = &BoltDB{}

// BoltDB implements the Database interface using a bolt database
type BoltDB struct {
	*bbolt.DB
	bucketName []byte
}

// NewBoltDB returns a new bolt database
func NewBoltDB(path string) (*BoltDB, error) {
	db, err := bbolt.Open(path, 0666, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating bolt db with path:%s", path)
	}

	return &BoltDB{DB: db, bucketName: []byte("default")}, nil
}

func (db *BoltDB) Get(keys []model.EncodedKey) (map[model.EncodedKey]model.TableValue, error) {
	res := make(map[model.EncodedKey]model.TableValue, len(keys))

	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(db.bucketName)

		for _, k := range keys {
			raw := b.Get(k.Bytes())
			if len(raw) == 0 {
				continue
			}
			val, err := model.ValueFromBytes(b.Get(k.Bytes()))
			if err != nil {
				return errors.Wrapf(err, "wrong record stored: %v", raw)
			}

			res[k] = val
		}

		return nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "an error occured when reading from bolt")
	}

	return res, nil
}

func (db *BoltDB) Set(batch map[model.EncodedKey]model.TableValue) error {
	err := db.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(db.bucketName)
		if err != nil {
			return errors.Wrapf(err, "error creating bucket")
		}

		for k, v := range batch {
			if err := b.Put(k.Bytes(), v.Bytes()); err != nil {
				return errors.Wrapf(err, "error setting (key: %v, val: %v)", k, v)
			}
		}
		return nil
	})

	if err != nil {
		return errors.Wrap(err, "an error occured when writing to bolt")
	}

	return nil
}
