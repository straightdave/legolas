package models

import (
	"github.com/fzzy/radix/redis"
)

type JobProxy struct {
	Queue string
	Rc    *redis.Client
}

// Pop up a job from the queue
func (j *JobProxy) Pop() (*Job, error) {
	data, err := j.Rc.Cmd("BLPOP", j.Queue, 0).List()
	if err != nil {
		return nil, err
	}
	return JobFromJson([]byte(data[1]))
}

// Append a job to the queue tail
func (j *JobProxy) Append(job *Job) error {
	data, err := job.Json()
	if err != nil {
		return err
	}
	return j.Rc.Cmd("RPUSH", j.Queue, string(data)).Err
}