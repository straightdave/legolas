package models

import (
	"encoding/json"
	"fmt"
)

type Job struct {
	CaseRunID string `json:"case_run_id"`
	Name      string `json:"name"`
	Snippet   string `json:"snippet"`
	PreAction string `json:"pre_action"`
}

// Unmarshal bytes into a job struct
func JobFromJson(data []byte) (job *Job, err error) {
	err = json.Unmarshal(data, job)
	return
}

// Marshal a job struct into bytes
func (job *Job) Json() ([]byte, error) {
	return json.Marshal(*job)
}

// Get job ID: <case_run_id>__<job_name>
func (job *Job) JobID() string {
	return fmt.Sprintf("%s__%s", job.CaseRunID, job.Name)
}

// Get previous job ID
func (job *Job) PrevJobID() string {
	return fmt.Sprintf("%s__%s", job.CaseRunID, job.PreAction)
}
