package ufx

import (
	"context"
	"flag"
	"go.uber.org/fx"
	"strings"
	"sync/atomic"
)

type ProbeParams struct {
	ReadinessCascade int64
}

func DecodeProbeParams(fset *flag.FlagSet) *ProbeParams {
	p := &ProbeParams{}
	fset.Int64Var(&p.ReadinessCascade, "probe.readiness.cascade", 5, "checker cascade")
	return p
}

type CheckerFunc func(ctx context.Context) error

func AsCheckerBuilder[T any](fn func(v T) (name string, cfn CheckerFunc)) any {
	return fx.Annotate(
		func(v T) named[CheckerFunc] {
			name, cfn := fn(v)
			return named[CheckerFunc]{Name: name, Val: cfn}
		},
		fx.ResultTags(`group:"ufx_checkers"`),
	)
}

// Probe is a check probe
type Probe interface {
	// CheckLiveness check liveness
	CheckLiveness() bool

	// CheckReadiness check readiness
	CheckReadiness(ctx context.Context) (s string, failed bool)
}

type probe struct {
	*ProbeParams

	checkers []named[CheckerFunc]
	failed   int64
}

type ProbeOptions struct {
	fx.In

	*ProbeParams

	Checkers []named[CheckerFunc] `group:"ufx_checkers"`
}

func NewProbe(opts ProbeOptions) Probe {
	return &probe{
		checkers:    opts.Checkers,
		ProbeParams: opts.ProbeParams,
	}
}

func (m *probe) CheckLiveness() bool {
	if m.ReadinessCascade > 0 {
		return m.failed < m.ReadinessCascade
	} else {
		return true
	}
}

func (m *probe) CheckReadiness(ctx context.Context) (result string, ready bool) {
	var results []string

	ready = true

	for _, checker := range m.checkers {
		var (
			name = checker.Name
			err  = checker.Val(ctx)
		)
		if err == nil {
			results = append(results, name+": OK")
		} else {
			results = append(results, name+": "+err.Error())
			ready = false
		}
	}

	result = strings.Join(results, "\n")

	if ready {
		atomic.StoreInt64(&m.failed, 0)
	} else {
		atomic.AddInt64(&m.failed, 1)
	}
	return
}
