/*
Test case run: running state of a test case
*/
package run

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"

	"legolas/common/config"
	"legolas/common/helpers"
	"legolas/common/models/action"
	"legolas/common/models/job"
	"legolas/common/models/testcase"
)

type Run struct {
	Id        string    `json:"id" bson:"id"`
	CasePath  string    `json:"case_path" bson:"case_path"`
	CaseName  string    `json:"case_name" bson:"case_name"`
	StartedAt time.Time `json:"started_at" bson:"started_at"`
	EndedAt   time.Time `json:"ended_at" bson:"ended_at"`
}

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

	// TODO: fix job proxy code
	jp := &job.JobProxy{
		Queue: config.Queue,
	}

	// push actions into queue
	for _, act := range acts {
		j := job.Job{
			CaseRunID:  runId,
			CasePath:   cpath,
			CaseName:   cname,
			ActionName: act.Name,
		}

		err = jp.Append(j)
		if err != nil {
			return
		}
	}

	return
}

func (r *Run) Save() (err error) {
	session, err := mgo.Dial(config.MongoHost)
	if err != nil {
		return
	}
	defer session.Close()

	// save run item first
	_, err = session.DB("legolas").C("runs").Upsert(bson.M{"id": r.Id, "case_path": r.CasePath, "case_name": r.CaseName}, *r)
	return
}

func GetRunsInOrder(cpath, cname string) (result []Run) {
	session, err := mgo.Dial(config.MongoHost)
	if err != nil {
		return
	}
	defer session.Close()

	err = session.DB("legolas").C("runs").Find(bson.M{"case_path": cpath, "case_name": cname}).Sort(...)
}
