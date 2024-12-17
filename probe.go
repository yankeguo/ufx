package ufx

import (
	"context"
	"strings"
	"sync"
	"sync/atomic"
)

// CheckerFunc is a check function for probe
type CheckerFunc func(ctx context.Context) error

// ProbeScope is a probe scope
type ProbeScope interface {
	// Check check probe
	Check(ctx context.Context) (result string, ready bool)

	// Add add checker to probe
	Add(name string, fn CheckerFunc)
}

// Probe is a check probe
type Probe interface {
	// Liveness check liveness
	Liveness() ProbeScope

	// Readiness check readiness
	Readiness() ProbeScope
}

type checkerItem struct {
	name string
	fn   CheckerFunc
}

type ProbeParams struct {
	Readiness struct {
		Cascade int `json:"cascade" default:"5" validate:"min=1"`
	} `json:"readiness"`
}

func ProbeParamsFromConf(conf Conf) (params ProbeParams, err error) {
	err = conf.Bind(&params, "probe")
	return
}

type probeScope struct {
	dep       *probeScope
	depTh     int64
	checkers  []checkerItem
	checkersL sync.Locker
	failed    int64
}

func newProbeScope(cascade *probeScope, threshold int64) *probeScope {
	return &probeScope{
		dep:       cascade,
		depTh:     threshold,
		checkersL: &sync.Mutex{},
	}
}

func (m *probeScope) Check(ctx context.Context) (result string, ready bool) {
	if m.dep != nil && m.depTh > 0 {
		if m.dep.failed >= m.depTh {
			return "Cascade failed", false
		}
	}

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

func (m *probeScope) Add(name string, fn CheckerFunc) {
	m.checkersL.Lock()
	defer m.checkersL.Unlock()

	m.checkers = append(m.checkers, checkerItem{name: name, fn: fn})
}

type probe struct {
	alive *probeScope
	ready *probeScope
}

func NewProbe(params ProbeParams) Probe {
	p := &probe{}
	p.ready = newProbeScope(nil, 0)
	p.alive = newProbeScope(p.ready, int64(params.Readiness.Cascade))
	return p
}

func (m *probe) Liveness() ProbeScope {
	return m.alive
}

func (m *probe) Readiness() ProbeScope {
	return m.ready
}
