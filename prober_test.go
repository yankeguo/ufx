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

	var m Prober

	fx.New(
		fx.Supply(r),
		fx.Supply(Conf(
			map[string]any{
				"prober": map[string]any{
					"readiness": map[string]any{
						"cascade": 2,
					},
				},
			},
		)),
		fx.Provide(
			ProberParamsFromConf,
			NewProber,
		),
		fx.Invoke(func(v Prober) {
			v.AddChecker("redis", func(ctx context.Context) error {
				if bad {
					return errors.New("test")
				}
				return nil
			})
		}),
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
