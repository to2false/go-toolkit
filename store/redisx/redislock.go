package redisx

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"strconv"
	"sync/atomic"
	"time"
)

// Copy from https://github.com/tal-tech/go-zero/blob/master/core/stores/redis/redislock.go
const (
	letters     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lockCommand = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    redis.call("SET", KEYS[1], ARGV[1], "PX", ARGV[2])
    return "OK"
else
    return redis.call("SET", KEYS[1], ARGV[1], "NX", "PX", ARGV[2])
end`
	delCommand = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end`
	randomLen       = 16
	tolerance       = 500 // milliseconds
	millisPerSecond = 1000
)

// A RedisLock is a redis lock.
type (
	RedisLock struct {
		store   *redis.Client
		seconds uint32
		key     string
		id      string
	}
	RedisLockOption func(lock *RedisLock)
)

// NewRedisLock returns a RedisLock.
func NewRedisLock(store *redis.Client, key string, options ...RedisLockOption) *RedisLock {
	r := &RedisLock{
		store: store,
		key:   key,
		id:    randomStr(randomLen),
	}

	for _, option := range options {
		option(r)
	}

	return r
}

func SetLockExpire(seconds uint32) RedisLockOption {
	return func(lock *RedisLock) {
		if lock == nil {
			return
		}
		lock.seconds = seconds
	}
}

// Acquire acquires the lock.
func (rl *RedisLock) Acquire(ctx context.Context) (bool, error) {
	seconds := atomic.LoadUint32(&rl.seconds)
	resp, err := rl.store.Eval(
		ctx,
		lockCommand, []string{rl.key}, []string{
			rl.id, strconv.Itoa(int(seconds)*millisPerSecond + tolerance),
		},
	).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("error on acquiring lock for %s, %s", rl.key, err.Error())
	} else if resp == nil {
		return false, nil
	}

	reply, ok := resp.(string)
	if ok && reply == "OK" {
		return true, nil
	}

	return false, nil
}

// Release releases the lock.
func (rl *RedisLock) Release(ctx context.Context) bool {
	resp, err := rl.store.Eval(ctx, delCommand, []string{rl.key}, []string{rl.id}).Result()

	if err != nil {
		return false
	}

	reply, ok := resp.(int64)
	if !ok {
		return false
	}

	return reply == 1
}

// UnsafeRelease Unsafe releases the lock ignore the value of lock
func (rl *RedisLock) UnsafeRelease(ctx context.Context) error {
	if err := rl.store.Del(ctx, rl.key).Err(); err != nil {
		return err
	}

	return nil
}

// SetExpire sets the expire.
func (rl *RedisLock) SetExpire(seconds int) {
	atomic.StoreUint32(&rl.seconds, uint32(seconds))
}

func randomStr(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint: gosec

	b := make([]byte, n)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))] // nolint
	}
	return string(b)
}
