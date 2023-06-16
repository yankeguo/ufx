package redisfx

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"ufx_redisfx",
	fx.Provide(
		DecodeParams,
		NewOptions,
		NewClient,
	),
	fx.Invoke(AddCheckerForClient),
)

var ModuleCluster = fx.Module(
	"ufx_redisfx_cluster",
	fx.Provide(
		DecodeClusterParams,
		NewClusterOptions,
		NewClusterClient,
	),
	fx.Invoke(AddCheckerForClusterClient),
)
