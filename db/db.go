package db

import "github.com/sfluor/musig/model"

// Database is an interface for storing fingerprint parts in a database and retrieving them
type Database interface {
	Get([]model.EncodedKey) (map[model.EncodedKey]model.TableValue, error)
	Set(map[model.EncodedKey]model.TableValue) error
	Close() error
}
