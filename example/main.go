package main

import (
	"github.com/yankeguo/ufx"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		ufx.Module,
		fx.Supply(
			ufx.Conf{},
		),
	).Run()
}
