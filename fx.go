package ufx

import "go.uber.org/fx"

var Module = fx.Module(
	"ufx",
	fx.Provide(
		ArgsFromCommandLine,
		NewFlagSet,
		BeforeParseFlagSet(DecodeProbeParams),
		BeforeParseFlagSet(DecodeRouterParams),
		BeforeParseFlagSet(DecodeServerParams),
		NewProbe,
		NewRouter,
		NewServer,
	),
	fx.Invoke(
		GuardedParseFlagSet(),
		SetupOTEL,
		Touch[Server],
	),
)
