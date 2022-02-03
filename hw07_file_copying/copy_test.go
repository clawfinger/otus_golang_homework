package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOffsetMoreThenLength(t *testing.T) {
	from := "testdata/input.txt"
	to := "testdata/out.txt"
	err := Copy(from, to, 100000, 0)
	require.Error(t, err)
	require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	os.Remove(to)
}

func TestOffsetEqualLength(t *testing.T) {
	from := "testdata/input.txt"
	to := "testdata/out.txt"
	err := Copy(from, to, 6617, 10)
	defer os.Remove(to)

	require.NoError(t, err)

	result, err := os.OpenFile(to, os.O_RDONLY, os.ModeAppend)

	require.NoError(t, err)

	fi, err := result.Stat()

	require.NoError(t, err)

	fileSize := fi.Size()
	require.Equal(t, int64(0), fileSize)
	result.Close()
}
