package models

import "gopkg.in/mgo.v2/bson"

type JobState struct {
	ActionRunId string            `json:"action_run_id"`
	State       string            `json:"state"`
	Results     map[string]string `json:"results"`
	Logs        []byte            `json:"logs"`
}

func (state *JobState) ToJson() ([]byte, error) {
	return bson.MarshalJSON(*state)
}

func JobStateFromJson(json []byte) (JobState, error) {
	var result JobState
	err := bson.UnmarshalJSON(json, &result)
	return result, err
}
