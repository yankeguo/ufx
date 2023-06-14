package ufx

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"testing"
)

func TestNewApp(t *testing.T) {
	type res struct {
		Hello string
	}
	r := &res{}
	var app Server

	a := fx.New(
		Module,
		fx.Supply(r),
		ReplaceArgs([]string{"--server.path.metrics", "/metrics"}),
		fx.Provide(
			AsCheckerBuilder(func(r *res) (name string, cfn CheckerFunc) {
				return "hello", func(ctx context.Context) error {
					return errors.New("bad")
				}
			}),
			AsRouteProvider(func(r *res) (name string, rfn HandlerFunc) {
				return "/hello", func(c Context) {
					c.Text("world")
				}
			}),
		),
		fx.Populate(&app),
	)

	require.NoError(t, a.Err())
}
