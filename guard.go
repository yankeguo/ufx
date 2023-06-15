package ufx

import "go.uber.org/fx"

// Guard is a type for ensuring invocation order
type Guard struct{}

func GuardParam[T any, U any](name string, fn func(t T) U) any {
	return fx.Annotate(
		func(_ []Guard, t T) U {
			return fn(t)
		},
		fx.ParamTags(`group:"`+name+`"`),
	)
}

func GuardParam21[I1 any, I2 any, O any](name string, fn func(i1 I1, i2 I2) O) any {
	return fx.Annotate(
		func(_ []Guard, t1 I1, t2 I2) O {
			return fn(t1, t2)
		},
		fx.ParamTags(`group:"`+name+`"`),
	)
}

func GuardResult[T any, U any](name string, fn func(t T) U) any {
	return fx.Annotate(
		func(t T) (Guard, U) {
			return Guard{}, fn(t)
		},
		fx.ResultTags(`group:"`+name+`"`),
	)
}
