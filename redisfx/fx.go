package redisfx

import (
	"github.com/guoyk93/ufx"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"ufx_redisfx",
	fx.Provide(
		ufx.AsFlagSetDecoder(DecodeParams),
		NewOptions,
		NewClient,
		ufx.AsCheckerBuilder(NewClientChecker),
	),
)

var ModuleCluster = fx.Module(
	"ufx_redisfx_cluster",
	fx.Provide(
		ufx.AsFlagSetDecoder(DecodeClusterParams),
		NewClusterOptions,
		NewClusterClient,
		ufx.AsCheckerBuilder(NewClusterClientChecker),
	),
)
