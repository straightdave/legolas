package run

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Run struct {
	Id         bson.ObjectId          `json:"_id" bson:"_id"`
	CaseId     bson.ObjectId          `json:"case_id" bson:"case_id"`
	StartedAt  time.Time              `json:"started_at" bson:"started_at"`
	EndedAt    time.Time              `json:"ended_at" bson:"ended_at"`
	Output     string                 `json:"output" bson:"output"`
	Context    map[string]interface{} `json:"context" bson:"context"`
	TracedData map[string]interface{} `json:"traced_data" bson:"traced_data"`
}

func (r *Run) Json() ([]byte, error) {
	return bson.MarshalJSON(*r)
}

func (r *Run) JsonPretty() ([]byte, error) {
	return json.MarshalIndent(*r, "", "    ")
}
