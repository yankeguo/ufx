package ufx

import "go.uber.org/fx"

var Module = fx.Module(
	"ufx",
	fx.Provide(
		LoadConf,
		NewProbeParamsFromConf,
		NewRouterParamsFromConf,
		NewServerParamsFromConf,
		NewProbe,
		NewRouter,
		NewServer,
	),
	fx.Invoke(SetupOTEL),
	fx.Invoke(func(Server) {}),
)
