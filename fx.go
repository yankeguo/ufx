package ufx

import "go.uber.org/fx"

var Module = fx.Module(
	"ufx",
	fx.Provide(
		ArgsFromCommandLine,
		NewFlagSet,
		AsFlagSetDecoder(DecodeProbeParams),
		AsFlagSetDecoder(DecodeRouterParams),
		NewProbe,
		NewRouter,
	),
	fx.Invoke(
		ParseFlagSet,
		SetupOTEL,
	),
)
