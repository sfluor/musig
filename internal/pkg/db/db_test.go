package db

import (
	"testing"

	"github.com/sfluor/musig/internal/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testTuple struct {
	key model.EncodedKey
	val model.TableValue
}

var testValues = []testTuple{
	{key: model.EncodedKey(1), val: model.TableValue{AnchorTimeMs: 1500, SongID: 100}},
	{key: model.EncodedKey(2), val: model.TableValue{AnchorTimeMs: 15, SongID: 42}},
	{key: model.EncodedKey(100), val: model.TableValue{AnchorTimeMs: 66, SongID: 999}},
	{key: model.EncodedKey(200), val: model.TableValue{AnchorTimeMs: 72, SongID: 100}},
	{key: model.EncodedKey(1000), val: model.TableValue{AnchorTimeMs: 1500, SongID: 999}},
	{key: model.EncodedKey(2000), val: model.TableValue{AnchorTimeMs: 65, SongID: 77}},
	{key: model.EncodedKey(30000), val: model.TableValue{AnchorTimeMs: 190, SongID: 100}},
	{key: model.EncodedKey(50000), val: model.TableValue{AnchorTimeMs: 38, SongID: 100}},
	{key: model.EncodedKey(428298445), val: model.TableValue{AnchorTimeMs: 3, SongID: 10}},
}

func testDatabase(t *testing.T, db Database) {
	// Close the database
	defer func() { require.NoError(t, db.Close()) }()

	t.Run("fingerprints", func(t *testing.T) {
		// Should return nothing without error
		res, err := db.Get(nil)
		require.NoError(t, err)
		assert.Len(t, res, 0)

		m1 := genTestMap(testValues[:4])
		err = db.Set(m1)
		require.NoError(t, err)

		keys := []model.EncodedKey{}
		for k := range m1 {
			keys = append(keys, k)
		}

		resMap, err := db.Get(keys)
		require.NoError(t, err)
		assert.Len(t, resMap, len(keys))
		assert.Equal(t, m1, resMap)

		m2 := genTestMap(testValues[4:])
		err = db.Set(m2)
		require.NoError(t, err)

		keys = []model.EncodedKey{}
		for k := range m2 {
			keys = append(keys, k)
		}

		resMap, err = db.Get(keys)
		require.NoError(t, err)
		assert.Len(t, resMap, len(keys))
		assert.Equal(t, m2, resMap)
	})

	t.Run("song_names", func(t *testing.T) {
		song1 := "my song !"
		song2 := "my second song !"

		name, err := db.GetSong(10)
		require.Error(t, err)
		assert.Empty(t, name)

		id, err := db.SetSong(song1)
		require.NoError(t, err)
		assert.True(t, id != 0)

		id2, err := db.SetSong(song2)
		require.NoError(t, err)
		assert.True(t, id2 != 0)

		name, err = db.GetSong(id2)
		require.NoError(t, err)
		assert.Equal(t, song2, name)
	})
}

func genTestMap(tuples []testTuple) map[model.EncodedKey]model.TableValue {
	m := make(map[model.EncodedKey]model.TableValue, len(tuples))
	for _, t := range tuples {
		m[t.key] = t.val
	}
	return m
}
