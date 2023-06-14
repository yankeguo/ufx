package halt

import (
	"errors"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestHalt(t *testing.T) {
	var err error
	func() {
		defer func() {
			err = recover().(error)
		}()

		String(
			"test",
			WithStatusCode(http.StatusTeapot),
			WithExtra("aaa", "bbb"),
			WithExtras(map[string]any{
				"ccc": "ddd",
				"eee": "fff",
			}),
		)
	}()
	m := ExtrasFromError(err)
	require.Equal(t, http.StatusTeapot, StatusCodeFromError(err))
	require.Equal(t, map[string]any{"message": "test", "aaa": "bbb", "ccc": "ddd", "eee": "fff"}, m)

	func() {
		defer func() {
			err = recover().(error)
		}()

		String(
			"test",
			WithBadRequest(),
			WithExtras(map[string]any{
				"ccc": "ddd",
				"eee": "fff",
			}),
			WithExtra("aaa", "bbb"),
			WithMessage("test2"),
		)
	}()
	m = ExtrasFromError(err)
	require.Equal(t, http.StatusBadRequest, StatusCodeFromError(err))
	require.Equal(t, map[string]any{"message": "test2", "aaa": "bbb", "ccc": "ddd", "eee": "fff"}, m)
}

func TestPanicError(t *testing.T) {
	var err error
	func() {
		defer func() {
			err = recover().(error)
		}()
		panic(errors.New("TEST1"))
	}()
	m := ExtrasFromError(err)
	require.Equal(t, http.StatusInternalServerError, StatusCodeFromError(err))
	require.Equal(t, map[string]any{"message": "TEST1"}, m)
}
