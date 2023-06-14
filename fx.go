package ufx

import "go.uber.org/fx"

var Module = fx.Module(
	"ufx",
	fx.Provide(
		ArgsFromCommandLine,
		NewFlagSet,
		AsFlagSetDecoder(DecodeProbeParams),
		AsFlagSetDecoder(DecodeRouterParams),
		AsFlagSetDecoder(DecodeServerParams),
		NewProbe,
		NewRouter,
		NewServer,
	),
	fx.Invoke(
		ParseFlagSet,
		SetupOTEL,
		touch[Server],
	),
)
