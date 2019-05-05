package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodingDecodingValue(t *testing.T) {
	for _, tv := range []TableValue{
		{123, 123},
		{1, 5},
		{9, 1000000},
		{3255, 42},
	} {
		encoded := tv.Bytes()
		arr, err := ValuesFromBytes(encoded)
		decoded := arr[0]
		require.NoError(t, err)
		encoded2 := decoded.Bytes()

		assert.Equal(t, tv, decoded)
		assert.Equal(t, encoded, encoded2)
	}
}

func TestDecodingMultipleValues(t *testing.T) {
	arr := []TableValue{
		{123, 123},
		{1, 5},
		{9, 1000000},
		{3255, 42},
	}

	var b []byte

	for _, v := range arr {
		b = append(b, v.Bytes()...)
	}

	res, err := ValuesFromBytes(b)
	require.NoError(t, err)
	assert.Equal(t, arr, res)
}
