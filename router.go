package ufx

import (
	"io/fs"
	"net/http"
	"path"
	"strings"
)

// HandlerFunc handler func with [Context] as argument
type HandlerFunc func(c Context)

// Router router interface
type Router interface {
	http.Handler

	HandleFS(root string, f fs.FS)

	HandleFunc(pattern string, fn HandlerFunc)
}

type router struct {
	RouterParams
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

type wrappedFS struct {
	root string
	f    fs.FS
}

func newWrappedFS(root string, f fs.FS) wrappedFS {
	root = strings.TrimSuffix(strings.TrimPrefix(root, "/"), "/")
	if root != "" {
		root = root + "/"
	}
	return wrappedFS{
		root: root,
		f:    f,
	}
}

func (w wrappedFS) Open(name string) (fs.File, error) {
	if w.root == "" {
		return w.f.Open(name)
	}
	if strings.HasPrefix(name, w.root) {
		return w.f.Open(name[len(w.root):])
	}
	return nil, fs.ErrNotExist
}

func (r *router) HandleFS(root string, f fs.FS) {
	if !strings.HasSuffix(root, "/") {
		root += "/"
	}
	if !strings.HasPrefix(root, "/") {
		root += "/"
	}

	var s http.Handler

	if root == "/" {
		s = http.FileServer(http.FS(f))
	} else {
		s = http.FileServer(http.FS(newWrappedFS(root, f)))
	}

	fs.WalkDir(f, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		r.m.Handle(path.Join(root, p), s)
		return nil
	})
}

func (r *router) HandleFunc(pattern string, fn HandlerFunc) {
	r.m.Handle(
		pattern,
		instrumentHTTPHandler(
			pattern,
			http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				c := newContext(rw, req)
				c.loggingResponse = r.Logging.Response
				c.loggingRequest = r.Logging.Request
				func() {
					defer c.Perform()
					fn(c)
				}()
			}),
		),
	)
}

type RouterParams struct {
	Concurrency int `json:"concurrency" default:"128" validate:"min=1"`
	Logging     struct {
		Response bool `json:"response"`
		Request  bool `json:"request"`
	} `json:"logging"`
}

func RouterParamsFromConf(conf Conf) (params RouterParams, err error) {
	err = conf.Bind(&params, "router")
	return
}

func NewRouter(opts RouterParams) Router {
	r := &router{
		RouterParams: opts,
		m:            &http.ServeMux{},
	}

	if opts.Concurrency > 0 {
		r.cc = make(chan struct{}, opts.Concurrency)
		for i := 0; i < opts.Concurrency; i++ {
			r.cc <- struct{}{}
		}
	}

	return r
}
