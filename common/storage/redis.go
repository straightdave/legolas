package storage

import (
	"fmt"
	"github.com/fzzy/radix/extra/pool"

	C "legolas/common/config"
)

var (
	redisPool *pool.Pool
)

func init() {
	var err error
	redisPool, err = pool.NewPool("tcp", C.RedisHost, C.RedisPoolSize)
	if err != nil {
		panic(fmt.Sprintf("Cannot create Redis pool at: %s\n", C.RedisHost))
	}
}

func GetRedisPool() *pool.Pool {
	return redisPool
}
