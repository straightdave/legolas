/*
   proxy of job state in redis
   job state in redis is a hash:
	key: job Id (case run id + action name)
	field "state" : state str (a string --- started/success/)
	field "error" : latest error message
	field "output" : output (a string --- stdour + stderr, only saved once when done)
	field "result" : result (json array string --- to store every type of data of python; only saved once when done)
	field "started_at" : string in redis, start time
	field "ended_at" : string, end time
*/
package models

import (
	"fmt"
	"github.com/fzzy/radix/redis"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

// for the purpose of saving to DB
type JobState struct {
	Id        string    `json:"uid" bson:"uid"`
	State     string    `json:"state" bson:"state"`
	Error     string    `json:"error" bson:"error"`
	Output    string    `json:"output" bson:"output"`
	Result    string    `json:"result" bson:"result"`
	StartedAt time.Time `json:"started_at" bson:"started_at"`
	EndedAt   time.Time `json:"ended_at" bson:"ended_at"`
}

type JobStateProxy struct {
	Id string
	Rc *redis.Client
}

func (s *JobStateProxy) State() (string, error) {
	return s.Rc.Cmd("HGET", s.Id, "state").Str()
}

func (s *JobStateProxy) MarkAs(s string) error {
	return s.Rc.Cmd("HSET", s.Id, "state", s).Err
}

// Mark a job as failed, along with error info particularly
func (s *JobStateProxy) Fail(reason string) error {
	s.Rc.Append("HSET", s.Id, "state", "failed")
	s.Rc.Append("HSET", s.Id, "error", reason)

	// call GetReply() to trigger batch operation
	// only return the first error if any (don't care about more)
	return s.Rc.GetReply().Err
}

func (s *JobStateProxy) SetOutput(content string) error {
	oriLog, err := s.Rc.Cmd("HGET", s.Id, "log").Str()
	if err != nil {
		L.Printf("cannot get original logs: %v\n", err)
		return
	}

	temp := strings.Join([]string{oriLog, newLog}, "\n------------\n")
	err = s.Rc.Cmd("HSET", s.Id, "log", temp).Err
	if err != nil {
		L.Printf("failed to add new logs: %v\n", err)
	}
}

func (s *JobStateProxy) GetOutput() string {
	oriLog, err := s.Rc.Cmd("HGET", s.Id, "log").Str()
	if err != nil {
		L.Printf("cannot get logs: %v\n", err)
		return ""
	}
	return oriLog
}

func (s *JobStateProxy) SaveToDB() error {

}
