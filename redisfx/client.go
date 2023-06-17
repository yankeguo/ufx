package redisfx

import (
	"context"
	"github.com/guoyk93/ufx"
	"github.com/redis/go-redis/v9"
)

type Params struct {
	URL string `json:"url" default:"redis://localhost:6379/0" validate:"required,url"`
}

func DecodeParams(conf ufx.Conf) (params Params, err error) {
	err = conf.Bind(&params, "redis")
	return
}

func NewOptions(params Params) (*redis.Options, error) {
	return redis.ParseURL(params.URL)
}

func NewClient(opts *redis.Options) *redis.Client {
	return redis.NewClient(opts)
}

func AddCheckerForClient(client *redis.Client, v ufx.Probe) {
	v.AddChecker("redis", func(ctx context.Context) error {
		return client.Ping(ctx).Err()
	})
}
