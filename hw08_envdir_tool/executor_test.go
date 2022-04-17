package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	cmd := []string{"echo"}
	code := RunCmd(cmd, nil)

	require.Equal(t, 0, code)
}
