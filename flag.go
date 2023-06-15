package ufx

import (
	"flag"
	"github.com/peterbourgon/ff/v3"
	"os"
)

const (
	guardParseFlagSet = "ufx_parse_flag_set"
)

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

// BeforeParseFlagSet wraps a flag set decoder function with joint points
func BeforeParseFlagSet[T any](fn func(fset *flag.FlagSet) *T) any {
	return GuardResult(guardParseFlagSet, fn)
}

// ParseFlagSet parses the flag set with ff
func ParseFlagSet(fset *flag.FlagSet, args Args) error {
	return ff.Parse(fset, args,
		ff.WithEnvVars(),
		ff.WithConfigFileFlag("conf"),
		ff.WithConfigFileParser(ff.PlainParser),
	)
}

func GuardedParseFlagSet() any {
	return GuardParam21(guardParseFlagSet, ParseFlagSet)
}
