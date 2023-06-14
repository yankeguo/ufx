package ufx

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"testing"
)

func TestAsCheckerBuilder(t *testing.T) {
	type res struct {
	}
	r := &res{}

	bad := true

	var m Probe

	fx.New(
		ReplaceArgs([]string{"-probe.readiness.cascade", "2"}),
		fx.Supply(r),
		fx.Provide(
			ArgsFromCommandLine,
			NewFlagSet,
			NewProbe,
			AsFlagSetDecoder(DecodeProbeParams),
			AsCheckerBuilder(func(v *res) (name string, cfn CheckerFunc) {
				return "redis", func(ctx context.Context) error {
					if bad {
						return errors.New("test")
					}
					return nil
				}
			}),
		),
		fx.Invoke(
			ParseFlagSet,
		),
		fx.Populate(&m),
	)

	require.True(t, m.CheckLiveness())

	s, failed := m.CheckReadiness(context.Background())
	require.True(t, failed)
	require.Equal(t, "redis: test", s)
	require.True(t, m.CheckLiveness())

	s, failed = m.CheckReadiness(context.Background())
	require.True(t, failed)
	require.Equal(t, "redis: test", s)
	require.False(t, m.CheckLiveness())

	bad = false

	s, failed = m.CheckReadiness(context.Background())
	require.False(t, failed)
	require.Equal(t, "redis: OK", s)
	require.True(t, m.CheckLiveness())
}
