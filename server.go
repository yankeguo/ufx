package ufx

import (
	"context"
	"flag"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/fx"
	"net/http"
	_ "net/http/pprof"
	"strings"
	"time"
)

// ServerParams params
type ServerParams struct {
	// Listen listen address
	Listen string

	// PathReadiness readiness check path
	PathReadiness string

	// PathLiveness liveness path
	PathLiveness string

	// PathMetrics metrics path
	PathMetrics string

	// DelayStart delay start
	DelayStart time.Duration

	// DelayStop delay stop
	DelayStop time.Duration
}

// DecodeServerParams create ServerParams from flag.FlagSet
func DecodeServerParams(fset *flag.FlagSet) (p *ServerParams) {
	p = &ServerParams{}
	fset.StringVar(&p.Listen, "server.listen", ":8080", "server listen address")
	fset.StringVar(&p.PathReadiness, "server.path.readiness", "/debug/ready", "server path readiness")
	fset.StringVar(&p.PathLiveness, "server.path.liveness", "/debug/alive", "server path liveness")
	fset.StringVar(&p.PathMetrics, "server.path.metrics", "/debug/metrics", "server path metrics")
	fset.DurationVar(&p.DelayStart, "server.delay.start", time.Second*3, "server delay start")
	fset.DurationVar(&p.DelayStop, "server.delay.stop", time.Second*3, "server delay stop")
	return p
}

// Server the main interface of [summer]
type Server interface {
	// Handler inherit [http.Handler]
	http.Handler
}

type server struct {
	*ServerParams

	Probe
	Router

	hProm http.Handler
}

func (a *server) serveReadiness(rw http.ResponseWriter, req *http.Request) {
	c := newContext(rw, req)
	defer c.Perform()

	s, failed := a.CheckReadiness(c)

	status := http.StatusOK
	if failed {
		status = http.StatusInternalServerError
	}

	c.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	c.Code(status)
	c.Text(s)
}

func (a *server) serveLiveness(rw http.ResponseWriter, req *http.Request) {
	c := newContext(rw, req)
	defer c.Perform()

	c.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	if a.CheckLiveness() {
		c.Code(http.StatusInternalServerError)
		c.Text("CASCADED FAILURE")
	} else {
		c.Code(http.StatusOK)
		c.Text("OK")
	}
}

func (a *server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// alive, ready, metrics
	if req.URL.Path == a.PathReadiness {
		// support readinessPath == livenessPath
		a.serveReadiness(rw, req)
		return
	} else if req.URL.Path == a.PathLiveness {
		a.serveLiveness(rw, req)
		return
	} else if req.URL.Path == a.PathMetrics {
		a.hProm.ServeHTTP(rw, req)
		return
	}

	// pprof
	if strings.HasPrefix(req.URL.Path, "/debug/pprof") {
		http.DefaultServeMux.ServeHTTP(rw, req)
		return
	}

	// serve with main handler
	a.Router.ServeHTTP(rw, req)
}

type ServerOptions struct {
	fx.In

	fx.Lifecycle

	*ServerParams

	Probe
	Router
}

// NewServer create an [Server] with [Option]
func NewServer(opts ServerOptions) Server {
	a := &server{
		ServerParams: opts.ServerParams,
		Probe:        opts.Probe,
		Router:       opts.Router,
		hProm:        promhttp.Handler(),
	}
	if opts.Lifecycle != nil {
		s := &http.Server{
			Addr:    opts.Listen,
			Handler: a,
		}
		opts.Lifecycle.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				chErr := make(chan error, 1)
				go func() {
					chErr <- s.ListenAndServe()
				}()
				select {
				case err := <-chErr:
					return err
				case <-ctx.Done():
					return s.Shutdown(ctx)
				case <-time.After(opts.DelayStart):
					return nil
				}
			},
			OnStop: func(ctx context.Context) error {
				time.Sleep(opts.DelayStop)
				return s.Shutdown(ctx)
			},
		})
	}
	return a
}
