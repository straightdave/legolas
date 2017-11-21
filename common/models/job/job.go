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
	"fmt"
)

type Job struct {
	CaseRunID  string `json:"case_run_id"`
	CasePath   string `json:"case_path"`
	CaseName   string `json:"case_name"`
	ActionName string `json:"action_name"`
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
