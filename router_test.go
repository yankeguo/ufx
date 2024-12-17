package ufx

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
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
			r.ServeMux().Handle("/", http.FileServer(http.Dir(filepath.Join("testdata", "router", "static"))))

			r.HandleFunc("/api/hello", func(c Context) {
				c.Text("world")
			})
		}),
		fx.Populate(&m),
	)
	require.NoError(t, a.Err())

	rw, req := httptest.NewRecorder(), httptest.NewRequest("GET", "/api/hello", nil)
	m.ServeHTTP(rw, req)

	require.Equal(t, "world", rw.Body.String())

	rw, req = httptest.NewRecorder(), httptest.NewRequest("GET", "/hello.txt", nil)
	m.ServeHTTP(rw, req)

	require.Equal(t, "hello world", rw.Body.String())

	rw, req = httptest.NewRecorder(), httptest.NewRequest("GET", "/hello/hello.txt", nil)
	m.ServeHTTP(rw, req)

	require.Equal(t, "hello world", rw.Body.String())

	rw, req = httptest.NewRecorder(), httptest.NewRequest("GET", "/hello/", nil)
	m.ServeHTTP(rw, req)

	t.Log(rw.Result().Status)
	t.Log(rw.Result().Header)

	require.Equal(t, "helloworld", strings.TrimSpace(rw.Body.String()))
}
