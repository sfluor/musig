package db

import "github.com/sfluor/musig/internal/pkg/model"

// Database is an interface for storing fingerprint parts in a database and retrieving them
type Database interface {
	Get([]model.EncodedKey) (map[model.EncodedKey][]model.TableValue, error)
	Set(map[model.EncodedKey]model.TableValue) error

	GetSongID(name string) (songID uint32, err error)
	GetSong(songID uint32) (name string, err error)
	SetSong(name string) (songID uint32, err error)

	Close() error
}
