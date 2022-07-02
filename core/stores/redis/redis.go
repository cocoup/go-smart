package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/breaker"
)

type Client interface {
	redis.Cmdable
}

type RedisConn interface {
}

type redisConn struct {
	client Client
	brk    breaker.Breaker
}

type clientMode func(conf Config) Client

func Single() clientMode {
	return func(conf Config) Client {
		return redis.NewClient(&redis.Options{})
	}
}

func Sentinel() clientMode {
	return func(conf Config) Client {
		return redis.NewFailoverClient(&redis.FailoverOptions{})
	}
}

func Cluster() clientMode {
	return func(conf Config) Client {
		return redis.NewClusterClient(&redis.ClusterOptions{})
	}
}

func NewConn(conf Config, client clientMode) RedisConn {
	return &redisConn{
		client: client(conf),
		brk:    breaker.NewBreaker(),
	}
}
