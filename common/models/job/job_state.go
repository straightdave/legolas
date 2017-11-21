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
package job

import (
	"time"
)

type JobState struct {
	CaseRunID  string    `json:"case_run_id"`
	CasePath   string    `json:"case_path"`
	CaseName   string    `json:"case_name"`
	ActionName string    `json:"action_name"`
	State      string    `json:"state" bson:"state"`
	Error      string    `json:"error" bson:"error"`
	Output     string    `json:"output" bson:"output"`
	Result     string    `json:"result" bson:"result"`
	StartedAt  time.Time `json:"started_at" bson:"started_at"`
	EndedAt    time.Time `json:"ended_at" bson:"ended_at"`
}
