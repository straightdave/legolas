/*
Job state represents the state of action run (job)
Job state is stored in mongodb.

    key: job Id (case run id + action name)
    field "state" : state str (a string --- started/success/)
    field "error" : latest error message
    field "output" : output (a string --- stdour + stderr, only saved once when done)
    field "result" : result (json array string --- to store every type of data of python; only saved once when done)
    field "started_at" : string in redis, start time
    field "ended_at" : string, end time
*/
package jobstate

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
)

type JobState struct {
	// keys to find a job state
	RunId    bson.ObjectId `json:"run_id" bson:"run_id"`
	ActionId bson.ObjectId `json:"action_id" bson:"action_id"`

	// other info
	State     string                 `json:"state" bson:"state"`
	Error     string                 `json:"error" bson:"error"`
	Output    string                 `json:"output" bson:"output"`
	Results   map[string]interface{} `json:"results" bson:"results"`
	StartedAt string                 `json:"started_at" bson:"started_at"`
	EndedAt   string                 `json:"ended_at" bson:"ended_at"`

	// informative for less queries from frontend
	// it stands for the name of the action at triggering time.
	// The action may change its name later, but this doesn't change.
	ActionName string `json:"action_name" bson:"action_name"`
}

func (js *JobState) JsonPretty() ([]byte, error) {
	return json.MarshalIndent(*js, "", "    ")
}
