package redisfx

import (
	"context"
	"flag"
	"github.com/guoyk93/ufx"
	"github.com/redis/go-redis/v9"
)

type ClusterParams struct {
	URL string
}

func DecodeClusterParams(fset *flag.FlagSet) *ClusterParams {
	opts := &ClusterParams{}
	fset.StringVar(&opts.URL, "redis.cluster.url", "redis://localhost:6379/0", "redis cluster url")
	return opts
}

func NewClusterOptions(p *ClusterParams) (*redis.ClusterOptions, error) {
	return redis.ParseClusterURL(p.URL)
}

func NewClusterClient(opts *redis.ClusterOptions) *redis.ClusterClient {
	return redis.NewClusterClient(opts)
}

func NewClusterClientChecker(client *redis.ClusterClient) (string, ufx.CheckerFunc) {
	return "redis", func(ctx context.Context) error {
		return client.Ping(ctx).Err()
	}
}
