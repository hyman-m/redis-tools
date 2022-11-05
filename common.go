package main

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient interface {
	Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd
	Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd
}
