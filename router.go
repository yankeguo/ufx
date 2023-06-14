package ufx

import (
	"flag"
	"go.uber.org/fx"
	"net/http"
)

// HandlerFunc handler func with [Context] as argument
type HandlerFunc func(c Context)

func AsRouteProvider[T any](fn func(v T) (pattern string, h HandlerFunc)) any {
	return fx.Annotate(
		func(v T) named[HandlerFunc] {
			pattern, rfn := fn(v)
			return named[HandlerFunc]{Name: pattern, Val: rfn}
		},
		fx.ResultTags(`group:"ufx_routes"`),
	)
}

type RouterParams struct {
	LoggingResponse bool
	LoggingRequest  bool
	Concurrency     int
}

func DecodeRouterParams(fset *flag.FlagSet) *RouterParams {
	p := &RouterParams{}
	fset.BoolVar(&p.LoggingRequest, "router.logging.request", false, "enable request logging")
	fset.BoolVar(&p.LoggingResponse, "router.logging.response", false, "enable response logging")
	fset.IntVar(&p.Concurrency, "router.concurrency", 128, "maximum concurrent requests")
	return p
}

type Router interface {
	http.Handler
}

type router struct {
	*RouterParams
	m  *http.ServeMux
	cc chan struct{}
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// concurrency
	if r.cc != nil {
		<-r.cc
		defer func() {
			r.cc <- struct{}{}
		}()
	}

	r.m.ServeHTTP(w, req)
}

type RouterOptions struct {
	fx.In

	*RouterParams

	Routes []named[HandlerFunc] `group:"ufx_routes"`
}

func NewRouter(opts RouterOptions) Router {
	r := &router{
		RouterParams: opts.RouterParams,
		m:            &http.ServeMux{},
	}

	if opts.Concurrency > 0 {
		r.cc = make(chan struct{}, opts.Concurrency)
		for i := 0; i < opts.Concurrency; i++ {
			r.cc <- struct{}{}
		}
	}

	for _, _item := range opts.Routes {
		item := _item
		r.m.Handle(
			item.Name,
			instrumentHTTPHandler(
				item.Name,
				http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
					c := newContext(rw, req)
					c.loggingResponse = opts.LoggingResponse
					c.loggingRequest = opts.LoggingRequest
					func() {
						defer c.Perform()
						item.Val(c)
					}()
				}),
			),
		)
	}

	return r
}
