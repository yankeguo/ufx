package ufx

import (
	"flag"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"os"
	"testing"
)

func TestNewFlagSet(t *testing.T) {
	s := NewFlagSet()
	require.NoError(t, s.Parse([]string{"--conf", "hello"}))
	f := s.Lookup("conf")
	require.Equal(t, "hello", f.Value.String())
}

func TestParseFlagSet(t *testing.T) {
	s := NewFlagSet()
	_ = s.String("ignore", "", "test")
	val := s.String("hello", "", "test")
	require.NoError(t, os.Setenv("HELLO", "WORLD"))
	require.NoError(t, ParseFlagSet(s, Args{"--ignore", "world"}))
	require.Equal(t, "WORLD", *val)
}

func TestAsFlagSetDecoder(t *testing.T) {
	type res struct {
		hello string
	}
	var r *res

	fx.New(
		fx.Supply(Args{"--hello", "world"}),
		fx.Provide(
			NewFlagSet,
			BeforeParseFlagSet(func(fset *flag.FlagSet) *res {
				r := &res{}
				fset.StringVar(&r.hello, "hello", "", "")
				return r
			}),
		),
		fx.Invoke(GuardedParseFlagSet()),
		fx.Populate(&r),
	)

	require.Equal(t, "world", r.hello)
}
