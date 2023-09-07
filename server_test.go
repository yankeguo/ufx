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
		fx.Supply(r, Conf{}),
		fx.Provide(
			ProberParamsFromConf,
			RouterParamsFromConf,
			ServerParamsFromConf,
			NewProber,
			NewRouter,
			NewServer,
		),
		fx.Invoke(func(r Router) {
			r.HandleFunc("/hello", func(c Context) {
				c.Text("world")
			})
		}),
		fx.Invoke(func(v Prober) {
			v.AddChecker("hello", func(ctx context.Context) error {
				return errors.New("bad")
			})
		}),
		fx.Populate(&app),
	)

	require.NoError(t, a.Err())
}
