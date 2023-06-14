package ufx

import "go.uber.org/fx"

var Module = fx.Module(
	"ufx",
	fx.Provide(
		ArgsFromCommandLine,
		NewFlagSet,
	),
	fx.Invoke(ParseFlagSet),
)
