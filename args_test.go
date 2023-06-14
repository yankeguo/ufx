package ufx

import (
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"testing"
)

func TestArgsFromCommandLine(t *testing.T) {
	var args Args
	fx.New(
		fx.Provide(ArgsFromCommandLine),
		fx.Populate(&args),
	)
	require.NotEmpty(t, args)
}

func TestReplaceArgs(t *testing.T) {
	var args Args
	fx.New(
		fx.Provide(ArgsFromCommandLine),
		ReplaceArgs([]string{"-hello", "world"}),
		fx.Populate(&args),
	)
	require.NotEmpty(t, args)
	require.Equal(t, Args{"-hello", "world"}, args)
}
