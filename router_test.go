package ufx

import (
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"net/http/httptest"
	"testing"
)

func TestNewRouter(t *testing.T) {
	type res struct {
	}
	r := &res{}

	var m Router

	a := fx.New(
		ReplaceArgs([]string{"-router.concurrency", "2"}),
		fx.Supply(r),
		fx.Provide(
			ArgsFromCommandLine,
			NewFlagSet,
			AsFlagSetDecoder(DecodeRouterParams),
			NewRouter,
			AsRouteProvider(func(r *res) (pattern string, h HandlerFunc) {
				return "/hello", func(c Context) {
					c.Text("world")
				}
			}),
		),
		fx.Invoke(
			ParseFlagSet,
		),
		fx.Populate(&m),
	)
	require.NoError(t, a.Err())

	rw, req := httptest.NewRecorder(), httptest.NewRequest("GET", "/hello", nil)
	m.ServeHTTP(rw, req)

	require.Equal(t, "world", rw.Body.String())
}
