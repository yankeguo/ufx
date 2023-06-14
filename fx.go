package ufx

import "go.uber.org/fx"

var Module = fx.Module(
	"ufx",
	fx.Provide(
		ArgsFromCommandLine,
		NewFlagSet,
		AsFlagSetDecoder(DecodeProbeParams),
		NewProbe,
	),
	fx.Invoke(
		ParseFlagSet,
		SetupOTEL,
	),
)
