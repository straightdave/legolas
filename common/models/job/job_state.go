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
package job

import (
	"encoding/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"

	"legolas/common/config"
)

type JobState struct {
	// keys to find a job state
	CaseRunID  string `json:"case_run_id" bson:"case_run_id"`
	ActionName string `json:"action_name" bson:"action_name"`

	// other info
	CasePath  string                 `json:"case_path" bson:"case_path"`
	CaseName  string                 `json:"case_name" bson:"case_name"`
	State     string                 `json:"state" bson:"state"`
	Error     string                 `json:"error" bson:"error"`
	Output    string                 `json:"output" bson:"output"`
	Results   map[string]interface{} `json:"results" bson:"results"` // set by python
	StartedAt time.Time              `json:"started_at" bson:"started_at"`
	EndedAt   time.Time              `json:"ended_at" bson:"ended_at"`
}

func NewJobState(caseRunId, actionName, cpath, cname string) JobState {
	return JobState{
		CaseRunID:  caseRunId,
		ActionName: actionName,
		CasePath:   cpath,
		CaseName:   cname,
		Results:    make(map[string]interface{}),
	}
}

func (js *JobState) JsonPretty() ([]byte, error) {
	return json.MarshalIndent(*js, "", "    ")
}

func (js *JobState) AddResult(name string, value interface{}) {
	js.Results[name] = value
}

func GetJobState(caseRunId, actionName string) (result JobState, err error) {
	session, err := mgo.Dial(config.MongoHost)
	if err != nil {
		return
	}
	defer session.Close()

	err = session.DB("legolas").C("jobstates").Find(bson.M{"case_run_id": caseRunId, "action_name": actionName}).One(&result)
	return
}

func (js *JobState) Save() (err error) {
	session, err := mgo.Dial(config.MongoHost)
	if err != nil {
		return
	}
	defer session.Close()

	_, err = session.DB("legolas").C("jobstates").Upsert(bson.M{"case_run_id": js.CaseRunID, "action_name": js.ActionName}, *js)
	return
}
