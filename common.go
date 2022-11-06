// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tools

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisClient .
type RedisClient interface {
	Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd
	Eval(ctx context.Context, script string, keys []string, args ...interface{}) *redis.Cmd
}
