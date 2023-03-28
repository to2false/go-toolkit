package redisx

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"github.com/to2false/go-toolkit/cache"
	"time"
)

type (
	RedisCache[T any] struct {
		node Node
	}
)

func (r RedisCache[T]) Del(ctx context.Context, keys ...cache.Key) error {
	tks := make([]string, 0, len(keys))
	for _, v := range keys {
		tks = append(tks, v.String())
	}

	return r.node.Client.Del(ctx, tks...).Err()
}

func (r RedisCache[T]) Get(ctx context.Context, key cache.Key) (result T, err error) {
	var res string
	res, err = r.node.Client.Get(ctx, key.String()).Result()
	if err != nil {
		if redis.Nil == err {
			err = nil
			return
		}

		return result, err
	}

	if err = json.Unmarshal([]byte(res), &result); err != nil {
		return result, err
	}

	return result, nil
}

func (r RedisCache[T]) Set(ctx context.Context, key cache.Key, val T, ttl time.Duration) error {
	res, err := json.Marshal(val)
	if err != nil {
		return err
	}

	return r.node.Client.Set(ctx, key.String(), res, ttl).Err()
}

func Cache[T any](node Node) cache.Cache[T] {
	return RedisCache[T]{
		node: node,
	}
}
