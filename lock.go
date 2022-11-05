package main

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

type RedisLock struct {
	Client     RedisClient
	Key        string
	uuid       string
	cancelFunc context.CancelFunc
}

func NewRedisLock(client RedisClient, key string) (*RedisLock, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	return &RedisLock{
		Client: client,
		Key:    key,
		uuid:   id.String(),
	}, nil
}

func (rl *RedisLock) TryLock(ctx context.Context) (bool, error) {
	succ, err := rl.Client.SetNX(ctx, rl.Key, rl.uuid, defaultExp).Result()
	if err != nil {
		return false, err
	}
	c, cancel := context.WithCancel(ctx)
	rl.cancelFunc = cancel
	rl.refresh(c)
	return succ, nil
}

func (rl *RedisLock) SpinLock(ctx context.Context, retryTimes int) (bool, error) {
	for i := 0; i < retryTimes; i++ {
		resp, err := rl.TryLock(ctx)
		if err != nil {
			return false, err
		}
		if resp {
			return resp, nil
		}
		time.Sleep(sleepDur)
	}
	return false, nil
}

func (rl *RedisLock) Unlock(ctx context.Context) (bool, error) {
	resp, err := NewTools(rl.Client).Cad(ctx, rl.Key, rl.uuid)
	if err != nil {
		return false, err
	}

	if resp {
		rl.cancelFunc()
	}
	return resp, nil
}

func (rl *RedisLock) refresh(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(defaultExp / 4)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				rl.Client.Expire(ctx, rl.Key, defaultExp)
			}
		}
	}()
}
