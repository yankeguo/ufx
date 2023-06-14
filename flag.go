package ufx

import (
	"flag"
	"github.com/peterbourgon/ff/v3"
	"go.uber.org/fx"
	"os"
)

type flagSetDecoderResult[T any] struct {
	fx.Out
	JP jointPoint `group:"ufx_flagset_decoder_jointpoints"`

	Value *T
}

// AsFlagSetDecoder wraps a flag set decoder function with joint points
func AsFlagSetDecoder[T any](fn func(fset *flag.FlagSet) *T) any {
	return func(fset *flag.FlagSet) flagSetDecoderResult[T] {
		return flagSetDecoderResult[T]{
			Value: fn(fset),
		}
	}
}

// NewFlagSet creates a new flag set
func NewFlagSet() *flag.FlagSet {
	name, _ := os.Executable()
	if name == "" {
		name = os.Args[0]
	}
	fset := flag.NewFlagSet(name, flag.ContinueOnError)
	_ = fset.String("conf", "", "config file (optional)")
	return fset
}

// ParseFlagSetOptions is the options for parsing flag set
type ParseFlagSetOptions struct {
	fx.In
	JP []jointPoint `group:"ufx_flagset_decoder_jointpoints"`

	FlagSet *flag.FlagSet
	Args    Args
}

// ParseFlagSet parses the flag set with ff
func ParseFlagSet(opts ParseFlagSetOptions) error {
	return ff.Parse(opts.FlagSet, opts.Args,
		ff.WithEnvVars(),
		ff.WithConfigFileFlag("conf"),
		ff.WithConfigFileParser(ff.PlainParser),
	)
}
