package gormfx

import (
	"github.com/guoyk93/ufx"
	"go.uber.org/fx"
)

var ModuleMySQL = fx.Module(
	"ufx_gormfx_mysql",
	fx.Provide(
		ufx.BeforeParseFlagSet(DecodeMySQLParams),
		NewMySQLDialector,
		NewConfig,
		NewClient,
		ufx.AsCheckerBuilder(NewClientChecker),
	),
)
