package redisx

import (
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

// todo @zengz mock redis client
func TestRedisLock(t *testing.T) {
	client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	key := "test"
	firstLock := NewRedisLock(client, key)
	firstLock.SetExpire(30)
	firstAcquire, err := firstLock.Acquire()
	assert.Nil(t, err)
	assert.True(t, firstAcquire)

	secondLock := NewRedisLock(client, key)
	secondLock.SetExpire(30)
	againAcquire, err := secondLock.Acquire()
	assert.Nil(t, err)
	assert.False(t, againAcquire)

	release := firstLock.Release()
	assert.True(t, release)

	endAcquire, err := secondLock.Acquire()
	assert.Nil(t, err)
	assert.True(t, endAcquire)
}

func TestNewRedisLock(t *testing.T) {
	type args struct {
		store   *redis.Client
		key     string
		options []RedisLockOption
	}
	tests := []struct {
		name string
		args args
		want *RedisLock
	}{
		{
			name: "test-second",
			args: args{
				store:   nil,
				key:     "test-three-second-expire",
				options: []RedisLockOption{SetLockExpire(3)},
			},
			want: &RedisLock{
				store:   nil,
				seconds: 3,
				key:     "test-three-second-expire",
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got := NewRedisLock(tt.args.store, tt.args.key, tt.args.options...)
				assert.Equal(t, got.seconds, tt.want.seconds)
				assert.Equal(t, got.key, tt.want.key)
				assert.Nil(t, got.store)
			},
		)
	}
}
