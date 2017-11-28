/*
Job: an action QUEUED and needs to be run

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
	"gopkg.in/mgo.v2/bson"
)

type Job struct {
	RunId        bson.ObjectId `json:"run_id" bson:"run_id"`
	ActionId     bson.ObjectId `json:"action_id" bson:"action_id"`
	PrevActionId bson.ObjectId `json:"prev_action_id" bson:"prev_action_id"`
}

func (job *Job) Json() ([]byte, error) {
	return json.Marshal(*job)
}

func (job *Job) JsonPretty() ([]byte, error) {
	return json.MarshalIndent(*job, "", "    ")
}

func FromJson(content []byte) (job Job, err error) {
	err = bson.UnmarshalJSON(content, &job)
	return
}

func (job *Job) Id() string {
	return fmt.Sprintf("%s__%s", job.RunId.Hex(), job.ActionId.Hex())
}
