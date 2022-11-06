// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tools

import (
	"context"
	"fmt"
	"time"
)

const (
	compareAndDeleteScript = `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
    	return redis.call("DEL", KEYS[1])
	else
    	return 0
	end
	`

	compareAndSwapScript = `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
    	return redis.call("SET", KEYS[1], ARGV[2])
	else
    	return 0
	end
	`

	compareAndSwapEXScript = `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
    	return redis.call("SET", KEYS[1], ARGV[2], %s ,ARGV[3])
	else
    	return 0
	end
	`

	compareAndSwapKeepTTLScript = `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
    	return redis.call("SET", KEYS[1], ARGV[2], "keepttl")
	else
    	return 0
	end
	`

	success = "OK"
)

// RedisTools .
type RedisTools struct {
	Client RedisClient
}

// NewTools create a new redis tools
func NewTools(client RedisClient) *RedisTools {
	return &RedisTools{Client: client}
}

// Cas compare and swap
func (r *RedisTools) Cas(ctx context.Context, key string, oldValue interface{},
	newValue interface{}) (bool, error) {

	res, err := r.Client.Eval(ctx, compareAndSwapScript, []string{key}, oldValue, newValue).Result()
	if err != nil {
		return false, err
	}
	if res == success {
		return true, nil
	}
	return false, nil
}

// CasEx compare and swap with timeout,
// If the timeout is 0, the timeout is not set,
// If the timeout is -1, keep timeout (redis >= 6.0).
func (r *RedisTools) CasEx(ctx context.Context, key string, oldValue interface{},
	newValue interface{}, expire time.Duration) (bool, error) {
	if expire == 0 {
		return r.Cas(ctx, key, oldValue, newValue)
	}

	var err error
	var res interface{}
	if usePrecise(expire) {
		res, err = r.Client.Eval(ctx, fmt.Sprintf(compareAndSwapEXScript, "PX"),
			[]string{key}, oldValue, newValue, formatMs(expire)).Result()
	} else if expire > 0 {
		res, err = r.Client.Eval(ctx, fmt.Sprintf(compareAndSwapEXScript, "EX"),
			[]string{key}, oldValue, newValue, formatSec(expire)).Result()
	} else {
		res, err = r.Client.Eval(ctx, compareAndSwapKeepTTLScript, []string{key},
			oldValue, newValue).Result()
	}

	if err != nil {
		return false, err
	}
	if res == success {
		return true, nil
	}
	return false, nil
}

// Cad compare and delete
func (r *RedisTools) Cad(ctx context.Context, key string, value interface{}) (bool, error) {
	res, err := r.Client.Eval(ctx, compareAndDeleteScript, []string{key}, value).Result()
	if err != nil {
		return false, err
	}
	if res == 0 {
		return false, nil
	}
	return true, nil
}

func usePrecise(dur time.Duration) bool {
	return dur < time.Second || dur%time.Second != 0
}

func formatMs(dur time.Duration) int64 {
	if dur > 0 && dur < time.Millisecond {
		return 1
	}
	return int64(dur / time.Millisecond)
}

func formatSec(dur time.Duration) int64 {
	if dur > 0 && dur < time.Second {
		return 1
	}
	return int64(dur / time.Second)
}
