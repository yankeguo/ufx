package redisfx

import (
	"context"
	"github.com/guoyk93/ufx"
	"github.com/redis/go-redis/v9"
)

type ClusterParams struct {
	URL string `json:"url"`
}

func DecodeClusterParams(conf ufx.Conf) (params ClusterParams, err error) {
	err = conf.Bind(&params, "redis", "cluster")
	return
}

func NewClusterOptions(params ClusterParams) (*redis.ClusterOptions, error) {
	return redis.ParseClusterURL(params.URL)
}

func NewClusterClient(opts *redis.ClusterOptions) *redis.ClusterClient {
	return redis.NewClusterClient(opts)
}

func AddCheckerForClusterClient(client *redis.ClusterClient, v ufx.Probe) {
	v.AddChecker("redis-cluster", func(ctx context.Context) error {
		return client.Ping(ctx).Err()
	})
}
