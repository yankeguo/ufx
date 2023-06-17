package main

import (
	"github.com/guoyk93/ufx"
	"github.com/guoyk93/ufx/redisfx"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

type app struct {
	r *redis.Client
}

func newApp(r *redis.Client) *app {
	return &app{r: r}
}

func addRoutesForApp(a *app, r ufx.Router) {
	r.HandleFunc("/hello", func(c ufx.Context) {
		c.Text("world")
	})
	r.HandleFunc("/get", a.routeGet)
	r.HandleFunc("/set", a.routeSet)
}

func (a *app) routeGet(c ufx.Context) {
	data := ufx.Bind[struct {
		Key string `json:"query_key"`
	}](c)
	c.Text(a.r.Get(c, data.Key).Val())
}

func (a *app) routeSet(c ufx.Context) {
	data := ufx.Bind[struct {
		Key string `json:"query_key"`
		Val string `json:"query_val"`
	}](c)
	c.Text(a.r.Set(c, data.Key, data.Val, 0).String())
}

func main() {
	fx.New(
		ufx.ProvideEmptyConf(),
		ufx.Module,
		redisfx.Module,
		fx.Provide(
			newApp,
		),
		fx.Invoke(addRoutesForApp),
	).Run()
}
