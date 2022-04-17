package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	env, err := ReadDir("./testdata/env")
	fmt.Println(env)
	require.Nil(t, err)

	v, ok := env["BAR"]
	require.True(t, ok)
	require.Equal(t, EnvValue{"bar", false}, v)

	v, ok = env["HELLO"]
	require.True(t, ok)
	require.Equal(t, EnvValue{"\"hello\"", false}, v)

	v, ok = env["EMPTY"]
	require.True(t, ok)
	require.Equal(t, EnvValue{"", false}, v)

	v, ok = env["UNSET"]
	require.True(t, ok)
	require.Equal(t, EnvValue{"", true}, v)
}

func TestDummydDir(t *testing.T) {
	_, err := ReadDir("missing")
	require.NotNil(t, err)
}
