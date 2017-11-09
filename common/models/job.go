package models

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
)

type Job struct {
	CaseRunID string `json:"case_run_id"`
	Name      string `json:"name"`
	Snippet   string `json:"snippet"`
	PreAction string `json:"pre_action"`
}

func (job *Job) ToJson() ([]byte, error) {
	return json.Marshal(*job)
}

func (job *Job) ToBson() ([]byte, error) {
	return bson.MarshalJSON(*job)
}

func JobFromJson(data []byte) (Job, error) {
	var job Job
	err := json.Unmarshal(data, &job)
	return job, err
}

func JobFromBson(data []byte) (Job, error) {
	var job Job
	err := bson.UnmarshalJSON(data, &job)
	return job, err
}
