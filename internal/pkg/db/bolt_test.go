package db

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBoltDB(t *testing.T) {
	testFile := "/tmp/test.db"
	db, err := NewBoltDB(testFile)
	require.NoError(t, err)

	testDatabase(t, db)

	err = os.Remove(testFile)
	require.NoError(t, err)
}
