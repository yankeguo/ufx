package ufx

import (
	"go.uber.org/fx"
	"os"
)

// Args is the command-line arguments
type Args []string

// ReplaceArgs override the command-line arguments for Fx
func ReplaceArgs(v []string) fx.Option {
	return fx.Replace(Args(v))
}

// ArgsFromCommandLine loads the flag set args from command-line arguments
func ArgsFromCommandLine() Args {
	return os.Args[1:]
}
