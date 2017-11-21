/*
Job: an action QUEUED and needs to be run.

When clients trigger one case run, it will:
1) generate case run id (unique)
2) iterate each action in the case
3) push each action into queue

The content is serialized in Json and stored in the queue.
From the content we could know:
1) case info
	- to get params of case level
2) action info
	- to get params of action level
	- to get template info
		- to get snippet

To use the Redis as little as possible, we use Redis for:
1) queue
2) temp storage for job running (in case job may access storage many times)
*/
package job

import (
	"encoding/json"
	"github.com/fzzy/radix/extra/pool"
	L "log"

	"legolas/common/config"
)

type Job struct {
	CaseRunID  string `json:"case_run_id"`
	CasePath   string `json:"case_path"`
	CaseName   string `json:"case_name"`
	ActionName string `json:"action_name"`
}

var redisPool *pool.Pool

func init() {
	var err error
	redisPool, err = pool.NewPool("tcp", config.RedisHost, config.RedisPoolSize)
	if err != nil {
		L.Fatalf("Cannot create Redis pool at %s: %v\n", config.RedisHost, err)
	}
}

func (job *Job) Json() ([]byte, error) {
	return json.Marshal(*job)
}

func (job *Job) JsonPretty() ([]byte, error) {
	return json.MarshalIndent(*job, "", "    ")
}

func FromJson(content []byte) (job Job, err error) {
	err = json.Unmarshal(content, &job)
	return
}

// Pop up a job from the queue
// Blocking if no job in the queue
func Pop() (job Job, err error) {
	rc, err := redisPool.Get()
	if err != nil {
		return
	}
	defer redisPool.Put(rc)

	data, err := rc.Cmd("BLPOP", config.Queue, 0).List()
	if err != nil {
		return
	}

	return FromJson([]byte(data[1]))
}

// Append a job to the queue's tail
func Append(job *Job) (err error) {
	rc, err := redisPool.Get()
	if err != nil {
		return
	}
	defer redisPool.Put(rc)

	data, err := job.Json()
	if err != nil {
		return
	}
	return rc.Cmd("RPUSH", config.Queue, string(data)).Err
}
