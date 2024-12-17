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
		fx.Supply(r),
		fx.Supply(Conf(
			map[string]any{
				"probe": map[string]any{
					"readiness": map[string]any{
						"cascade": 2,
					},
				},
			},
		)),
		fx.Provide(
			ProbeParamsFromConf,
			NewProbe,
		),
		fx.Invoke(func(v Probe) {
			v.Readiness().Add("redis", func(ctx context.Context) error {
				if bad {
					return errors.New("test")
				}
				return nil
			})
		}),
		fx.Populate(&m),
	)

	ctx := context.Background()

	_, ok := m.Liveness().Check(ctx)

	require.True(t, ok)

	s, ready := m.Readiness().Check(context.Background())
	require.False(t, ready)
	require.Equal(t, "redis: test", s)

	_, ok = m.Liveness().Check(ctx)
	require.True(t, ok)

	s, ready = m.Readiness().Check(context.Background())
	require.False(t, ready)
	require.Equal(t, "redis: test", s)

	_, ok = m.Liveness().Check(ctx)
	require.False(t, ok)

	bad = false

	s, ready = m.Readiness().Check(context.Background())
	require.True(t, ready)
	require.Equal(t, "redis: OK", s)

	_, ok = m.Liveness().Check(ctx)
	require.True(t, ok)
}
