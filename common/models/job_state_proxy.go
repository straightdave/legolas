package models

import (
	"github.com/fzzy/radix/redis"
	L "log"
	"strings"
)

/*
   proxy of job state in redis
*/

type JobStateProxy struct {
	Id string
	Rc *redis.Client
}

// Return the state string of a job
func (s *JobStateProxy) State() string {
	sta, err := s.Rc.Cmd("HGET", s.Id, "state").Str()
	if err != nil {
		L.Printf("cannot get state of job [%s]: %v\n", s.Id, err)
		return "err"
	}
	return sta
}

// Mark a job as In-Progress
func (s *JobStateProxy) InProgress() error {
	return s.Rc.Cmd("HSET", s.Id, "state", "inprogress").Err
}

// Mark a job as Done
func (s *JobStateProxy) Done() {
	err := s.Rc.Cmd("HSET", s.Id, "state", "done").Err
	if err != nil {
		L.Printf("failed to set job run [%s] state as done: %v\n", s.Id, err)
	}
}

// Mark a job as Failed, with error info
func (s *JobStateProxy) Fail(reason string) {
	s.Rc.Append("HSET", s.Id, "state", "failed")
	s.Rc.Append("HSET", s.Id, "error", reason)

	// only care about the first setting operation
	err := s.Rc.GetReply().Err
	if err != nil {
		L.Printf("failed to set job [%s] state to fail: %v\n", s.Id, err)
	}
}

// Append log to an existing job
func (s *JobStateProxy) AppendLog(newLog string) {
	oriLog, err := s.Rc.Cmd("HGET", s.Id, "log").Str()
	if err != nil {
		L.Printf("cannot get original logs: %v\n", err)
		return
	}

	temp := strings.Join([]string{oriLog, newLog}, "\n")
	err = s.Rc.Cmd("HSET", s.Id, "log", temp).Err
	if err != nil {
		L.Printf("failed to add new logs: %v\n", err)
	}
}

/*
   proxy of case run state in redis
*/

type CaseRunStateProxy struct {
	Id string
	Rc *redis.Client
}

// TODO
