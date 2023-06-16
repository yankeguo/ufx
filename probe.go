package ufx

import (
	"context"
	"strings"
	"sync/atomic"
)

type CheckerFunc func(ctx context.Context) error

// Probe is a check probe
type Probe interface {
	// CheckLiveness check liveness
	CheckLiveness() bool

	// CheckReadiness check readiness
	CheckReadiness(ctx context.Context) (s string, failed bool)

	// AddChecker add checker
	AddChecker(name string, fn CheckerFunc)
}

type probeItem struct {
	name string
	fn   CheckerFunc
}

type ProbeParams struct {
	Readiness struct {
		Cascade int `yaml:"cascade" default:"5" validate:"min=1"`
	} `yaml:"readiness"`
}

func NewProbeParamsFromConf(conf Conf) (params ProbeParams, err error) {
	err = conf.Bind(&params, "probe")
	return
}

type probe struct {
	ProbeParams

	checkers []probeItem
	failed   int64
}

func NewProbe(params ProbeParams) Probe {
	return &probe{
		ProbeParams: params,
	}
}

func (m *probe) CheckLiveness() bool {
	if m.Readiness.Cascade > 0 {
		return m.failed < int64(m.Readiness.Cascade)
	} else {
		return true
	}
}

func (m *probe) CheckReadiness(ctx context.Context) (result string, ready bool) {
	var results []string

	ready = true

	for _, checker := range m.checkers {
		var (
			name = checker.name
			err  = checker.fn(ctx)
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

func (m *probe) AddChecker(name string, fn CheckerFunc) {
	m.checkers = append(m.checkers, probeItem{name: name, fn: fn})
}
