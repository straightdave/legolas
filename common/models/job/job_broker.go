package job

import (
	"github.com/fzzy/radix/extra/pool"
	"legolas/common/config"
)

var redisPool *pool.Pool

func SetRedisPool(p *pool.Pool) {
	redisPool = p
}

func View() (result []Job, err error) {
	rc, err := redisPool.Get()
	if err != nil {
		return
	}
	defer redisPool.CarefullyPut(rc, &err)

	data, err := rc.Cmd("LRANGE", config.Queue, 0, -1).List()
	if err != nil {
		return
	}

	result = []Job{}
	for _, item := range data {
		j, _ := FromJson([]byte(item))
		result = append(result, j)
	}
	return
}

func Pop() (job Job, err error) {
	rc, err := redisPool.Get()
	if err != nil {
		return
	}
	defer redisPool.CarefullyPut(rc, &err)

	data, err := rc.Cmd("BLPOP", config.Queue, 0).List()
	if err != nil {
		return
	}
	return FromJson([]byte(data[1])) // data[0] is col name
}

func Append(job *Job) (err error) {
	rc, err := redisPool.Get()
	if err != nil {
		return
	}
	defer redisPool.CarefullyPut(rc, &err)

	data, err := job.Json()
	if err != nil {
		return
	}
	return rc.Cmd("RPUSH", config.Queue, string(data)).Err
}
