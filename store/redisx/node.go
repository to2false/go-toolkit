package redisx

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/to2false/go-toolkit/cache"
	"strings"
	"time"
)

type (
	Config struct {
		Addr     string `json:"addr"`
		DB       int    `json:"db"`
		Password string `json:"password"`

		PoolSize   int `json:"pool_size"`
		MaxRetries int `json:"max_retries"`
	}

	Node struct {
		*redis.Client
	}
	Nodes []Node
)

func MustNew(ctx context.Context, c *Config) Nodes {
	nodes, err := New(ctx, c)
	if err != nil {
		panic(err)
	}

	return nodes
}

func New(ctx context.Context, c *Config) (Nodes, error) {
	addrs := strings.Split(c.Addr, ",")
	if len(addrs) == 0 {
		return nil, fmt.Errorf("empty redis addr")
	}

	nodes := make(Nodes, 0, len(addrs))
	for _, addr := range addrs {
		if addr == "" {
			continue
		}
		cli := redis.NewClient(
			&redis.Options{
				Addr:     addr,
				DB:       c.DB,
				Password: c.Password,

				MaxRetries: c.MaxRetries,
				PoolSize:   c.PoolSize,
			},
		)

		ping := cli.Ping(ctx)
		if ping.Val() == "" {
			return nil, fmt.Errorf("redis connect addr %s err: %s", addr, ping.Err())
		}

		nodes = append(nodes, Node{cli})
	}

	return nodes, nil
}

func (nodes Nodes) Get() Node {
	return RoundRobin(nodes)
}

func (node Node) Locker(key cache.Key, ttl time.Duration) *RedisLock {
	return NewRedisLock(node.Client, key.String(), SetLockExpire(uint32(ttl.Seconds())))
}
