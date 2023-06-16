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
		fx.Supply(r, Conf{}),
		fx.Provide(
			NewRouterParamsFromConf,
			NewRouter,
		),
		fx.Invoke(func(r Router) {
			r.HandleFunc("/hello", func(c Context) {
				c.Text("world")
			})
		}),
		fx.Populate(&m),
	)
	require.NoError(t, a.Err())

	rw, req := httptest.NewRecorder(), httptest.NewRequest("GET", "/hello", nil)
	m.ServeHTTP(rw, req)

	require.Equal(t, "world", rw.Body.String())
}
