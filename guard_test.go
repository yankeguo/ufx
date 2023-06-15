package ufx

import (
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"testing"
)

func TestGuard(t *testing.T) {
	type O struct{}
	type I struct{}
	type A struct{}
	type B struct{}
	type C struct{}

	created := make([]bool, 3)

	fx.New(
		fx.Supply(O{}),
		fx.Provide(
			GuardResult("test", func(O) A {
				created[0] = true
				return A{}
			}),
			GuardResult("test", func(O) B {
				created[1] = true
				return B{}
			}),
			GuardResult("test", func(O) C {
				created[2] = true
				return C{}
			}),
			GuardParam("test", func(O) I {
				return I{}
			}),
		),
		fx.Invoke(func(I) {}),
	)

	require.True(t, created[0])
	require.True(t, created[1])
	require.True(t, created[2])
}
