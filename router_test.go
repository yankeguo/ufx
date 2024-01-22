package ufx

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

func TestNewRouter(t *testing.T) {
	type res struct {
	}
	r := &res{}

	var m Router

	a := fx.New(
		fx.Supply(r, Conf{}),
		fx.Provide(
			RouterParamsFromConf,
			NewRouter,
		),
		fx.Invoke(func(r Router) {
			r.HandleFS("/static", os.DirFS(filepath.Join("testdata", "router", "static")))

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

	rw, req = httptest.NewRecorder(), httptest.NewRequest("GET", "/static/hello.txt", nil)
	m.ServeHTTP(rw, req)

	require.Equal(t, "hello world", rw.Body.String())

	rw, req = httptest.NewRecorder(), httptest.NewRequest("GET", "/static/hello/hello.txt", nil)
	m.ServeHTTP(rw, req)

	require.Equal(t, "hello world", rw.Body.String())
}
