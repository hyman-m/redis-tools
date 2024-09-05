# redis-tools

[![Go Report Card](https://goreportcard.com/badge/github.com/wanzoma/redis-tools)](https://goreportcard.com/report/github.com/wanzoma/redis-tools)&nbsp;![GitHub top language](https://img.shields.io/github/languages/top/wanzoma/redis-tools)&nbsp;![GitHub](https://img.shields.io/github/license/wanzoma/balancer)&nbsp;[![CodeFactor](https://www.codefactor.io/repository/github/wanzoma/balancer/badge)](https://www.codefactor.io/repository/github/wanzoma/redis-tools)&nbsp;![go_version](https://img.shields.io/badge/go%20version-1.19-yellow)

redis-tools is a collection of redis tools, including `distributed lock`, `cas`, `casEx`, `cad` .

# Quick Start
Fisrt, create a demo and import the redis-tools and redis client :  

```shell
> go mod init demo

> go get github.com/wanzoma/redis-tools
> go get github.com/go-redis/redis/v8
```

## Distributed lock
The `trylock` case :
```go
package main

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	tools "github.com/zehuamama/redis-tools"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	disLock, err := tools.NewRedisLock(client, "lock resource")
	if err != nil {
		log.Fatal(err)
	}

	succ, err := disLock.TryLock(context.Background())
	if err != nil {
		log.Println(err)
        return
	}

	if succ {
		defer disLock.Unlock(context.Background())
	}
}

```
and `spinlock` case ：
```go
    succ, err := disLock.SpinLock(context.Background(), 5)  // retry 5 times
	if err != nil {
		log.Println(err)
        return
	}

	if succ {
		defer disLock.Unlock(context.Background())
	}
```

## Redis Tools
`compare and swap` case :
```go
func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	succ, err := tools.NewTools(client).Cas(context.Background(), "cas_key", "old value", "new value")
	if err != nil {
		log.Println(err)
		return
	}

    ...
}
```

and `compare and delete` case :
```go
    succ, err := tools.NewTools(client).Cad(context.Background(), "cas_key", "old value")
	if err != nil {
		log.Println(err)
		return
	}
```

## Contributing

If you are intersted in contributing to redis-tools, please see here: [CONTRIBUTING](https://github.com/zehuamama/redis-tools/blob/main/CONTRIBUTING.md)

## License

redis-tools is licensed under the term of the [BSD 2-Clause License](https://github.com/zehuamama/redis-tools/blob/main/LICENSE)
