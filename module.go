package ufx

import "go.uber.org/fx"

var Module = fx.Module(
	"ufx",
	fx.Provide(
		ProbeParamsFromConf,
		RouterParamsFromConf,
		ServerParamsFromConf,
		NewProbe,
		NewRouter,
		NewServer,
	),
	fx.Invoke(SetupTelemetry),
	fx.Invoke(func(Server) {}),
)
