package main

import (
	"context"
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
    	return redis.call("SET", KEYS[1], ARGV[2], "EX" ,ARGV[3])
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

	compareAndSwapPXScript = `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
    	return redis.call("SET", KEYS[1], ARGV[2], "PX" ,ARGV[3])
	else
    	return 0
	end
	`

	success = "OK"
)

type RedisTools struct {
	Client RedisClient
}

func NewTools(client RedisClient) *RedisTools {
	return &RedisTools{Client: client}
}

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

func (r *RedisTools) CasEx(ctx context.Context, key string, oldValue interface{},
	newValue interface{}, expire time.Duration) (bool, error) {
	if expire == 0 {
		return r.Cas(ctx, key, oldValue, newValue)
	}

	var err error
	var res interface{}
	if usePrecise(expire) {
		res, err = r.Client.Eval(ctx, compareAndSwapPXScript, []string{key},
			oldValue, newValue, formatMs(expire)).Result()
	} else if expire > 0 {
		res, err = r.Client.Eval(ctx, compareAndSwapEXScript, []string{key},
			oldValue, newValue, formatSec(expire)).Result()
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
