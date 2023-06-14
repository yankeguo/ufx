package ufx

import (
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"testing"
)

func TestModule(t *testing.T) {
	app := fx.New(Module, ReplaceArgs([]string{}))
	require.NoError(t, app.Err())
}
