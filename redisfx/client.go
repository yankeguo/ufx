package redisfx

import (
	"context"
	"flag"
	"github.com/guoyk93/ufx"
	"github.com/redis/go-redis/v9"
)

type Params struct {
	URL string
}

func DecodeParams(fset *flag.FlagSet) *Params {
	opts := &Params{}
	fset.StringVar(&opts.URL, "redis.url", "redis://localhost:6379/0", "redis url")
	return opts
}

func NewOptions(p *Params) (*redis.Options, error) {
	return redis.ParseURL(p.URL)
}

func NewClient(opts *redis.Options) *redis.Client {
	return redis.NewClient(opts)
}

func NewClientChecker(client *redis.Client) (string, ufx.CheckerFunc) {
	return "redis", func(ctx context.Context) error {
		return client.Ping(ctx).Err()
	}
}
