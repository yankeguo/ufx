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
		ReplaceArgs("-probe.readiness.cascade", "2"),
		fx.Supply(r),
		fx.Provide(
			ArgsFromCommandLine,
			NewFlagSet,
			NewProbe,
			BeforeParseFlagSet(DecodeProbeParams),
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
			GuardedParseFlagSet(),
		),
		fx.Populate(&m),
	)

	require.True(t, m.CheckLiveness())

	s, ready := m.CheckReadiness(context.Background())
	require.False(t, ready)
	require.Equal(t, "redis: test", s)
	require.True(t, m.CheckLiveness())

	s, ready = m.CheckReadiness(context.Background())
	require.False(t, ready)
	require.Equal(t, "redis: test", s)
	require.False(t, m.CheckLiveness())

	bad = false

	s, ready = m.CheckReadiness(context.Background())
	require.True(t, ready)
	require.Equal(t, "redis: OK", s)
	require.True(t, m.CheckLiveness())
}
