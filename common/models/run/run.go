/*
Test case run: running info of a test case
*/
package run

import (
	"encoding/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"

	"legolas/common/config"
	"legolas/common/helpers"
	"legolas/common/models/action"
	"legolas/common/models/job"
)

type Run struct {
	Id        string    `json:"case_run_id" bson:"case_run_id"`
	CasePath  string    `json:"case_path" bson:"case_path"`
	CaseName  string    `json:"case_name" bson:"case_name"`
	StartedAt time.Time `json:"started_at" bson:"started_at"`
	EndedAt   time.Time `json:"ended_at" bson:"ended_at"`
}

func (r *Run) JsonPretty() ([]byte, error) {
	return json.MarshalIndent(*r, "", "    ")
}

func GetRuns(cpath, cname string) (result []Run, err error) {
	session, err := mgo.Dial(config.MongoHost)
	if err != nil {
		return
	}
	defer session.Close()

	err = session.DB("legolas").C("runs").Find(bson.M{"case_path": cpath, "case_name": cname}).All(&result)
	return
}

func GetRun(id string) (result Run, err error) {
	session, err := mgo.Dial(config.MongoHost)
	if err != nil {
		return
	}
	defer session.Close()

	err = session.DB("legolas").C("runs").Find(bson.M{"case_run_id": id}).One(&result)
	return
}

// create a test case run:
// push all actions into job queue
func NewRun(cpath, cname string) (r Run, err error) {
	r = Run{
		Id:        helpers.RandStringRunes(config.RunIdLength),
		CasePath:  cpath,
		CaseName:  cname,
		StartedAt: time.Now(),
	}

	session, err := mgo.Dial(config.MongoHost)
	if err != nil {
		return
	}
	defer session.Close()

	// save run item first
	err = session.DB("legolas").C("runs").Insert(r)
	if err != nil {
		return
	}

	// get all actions
	acts, err := action.GetActions(cpath, cname)
	if err != nil {
		return
	}

	// push actions into queue
	// Order matters!
	prevAct := ""
	for _, act := range acts {
		j := &job.Job{
			CaseRunID:  r.Id,
			CasePath:   cpath,
			CaseName:   cname,
			ActionName: act.Name,
			PrevAction: prevAct,
		}

		err = job.Append(j)
		if err != nil {
			return
		}

		prevAct = act.Name
	}
	return
}

// update in mongo. no effect on redis queue
func (r *Run) Save() (err error) {
	session, err := mgo.Dial(config.MongoHost)
	if err != nil {
		return
	}
	defer session.Close()

	_, err = session.DB("legolas").C("runs").Upsert(bson.M{"case_run_id": r.Id, "case_path": r.CasePath, "case_name": r.CaseName}, *r)
	return
}
