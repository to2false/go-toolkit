package cache

import (
	"context"
	"time"
)

type (
	Cache[T any] interface {
		// Del deletes cached values with keys.
		Del(ctx context.Context, keys ...Key) error
		// Get gets the cache with key and fills into v.
		Get(ctx context.Context, key Key) (T, error)
		// Set sets the cache with key and v, using c.expiry.
		Set(ctx context.Context, key Key, val T, ttl time.Duration) error
	}
)
