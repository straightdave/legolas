package models

import "gopkg.in/mgo.v2/bson"

type CaseRunState struct {
	CaseRunId string                 `json:"case_run_id"`
	Context   map[string]interface{} `json:"context"`
}

func (state *CaseRunState) ToJson() ([]byte, error) {
	return bson.MarshalJSON(*state)
}

func CaseRunStateFromJson(json []byte) (CaseRunState, error) {
	var result CaseRunState
	err := bson.UnmarshalJSON(json, &result)
	return result, err
}
