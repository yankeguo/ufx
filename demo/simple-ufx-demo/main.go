package main

import (
	"context"
	"github.com/guoyk93/ufx"
	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

func createRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{})
}

func checkRedisClient(c *redis.Client) (string, ufx.CheckerFunc) {
	return "redis", func(ctx context.Context) error {
		return c.Ping(ctx).Err()
	}
}

type app struct {
	r *redis.Client
}

func newApp(r *redis.Client) *app {
	return &app{r: r}
}

func appRouteGet(a *app) (string, ufx.HandlerFunc) {
	return "/get", func(c ufx.Context) {
		data := ufx.Bind[struct {
			Key string `json:"query_key"`
		}](c)
		c.Text(a.r.Get(c, data.Key).Val())
	}
}

func appRouteSet(a *app) (string, ufx.HandlerFunc) {
	return "/set", func(c ufx.Context) {
		data := ufx.Bind[struct {
			Key string `json:"query_key"`
			Val string `json:"query_val"`
		}](c)
		c.Text(a.r.Set(c, data.Key, data.Val, 0).String())
	}
}

func main() {
	fx.New(
		ufx.Module,
		fx.Provide(
			createRedisClient,
			ufx.AsCheckerBuilder(checkRedisClient),
			newApp,
			ufx.AsRouteProvider(appRouteGet),
			ufx.AsRouteProvider(appRouteSet),
		),
	).Run()
}
