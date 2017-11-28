/*
   CaseRun: data to store at DB, representing a test case run object
   it's not in the Redis queue where stores only action jobs.
   however it goes to Redis as simple key-value hash for Legolas to track status.
*/
package models

import (
	"fmt"
	"github.com/fzzy/radix/extra/pool"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	CM "legolas/common/models"
	L "log"
	"time"
)

const (
	redisHost = "localhost:6379"
	queueName = "actionjobs"
)

var (
	poolSize  = 10
	redisPool *pool.Pool
)

func init() {
	var err error
	redisPool, err = pool.NewPool("tcp", redisHost, poolSize)
	if err != nil {
		L.Printf("Cannot create Redis pool at: %s\n", redisHost)
	}
}

type Status struct {
	Name string    `json:"name"`
	At   time.Time `json:"at"`
	Msg  string    `json:"message"`
}

type CaseRun struct {
	Id          string   `json:"uid" bson:"uid"`
	Path        string   `json:"path" bson:"path"`
	Name        string   `json:"name" bson:"name"`
	Stati       []Status `json:"stati" bson:"stati"`
	ActionNames []string `json:"action_names" bson:"action_names"`
	actions     []Action // private, for storing action data
}

// Key function!
// create new instance as well as saving it to DB and Redis
// more important, saving all actions into job queue in Redis
func NewCaseRun(c *Case) (cr CaseRun, err error) {
	var timeStamp = time.Now().Format("20060102150405")
	var suffix = C.RandSuffix4()

	cr.Id = fmt.Sprintf("%s_%s", timeStamp, suffix) // unique id (hopefully)
	cr.Path = c.Path
	cr.Name = c.Name
	cr.NewStatus("creating", "")
	acts, err := FindActions(c.Name, c.Path)
	if err != nil {
		cr.NewStatus("creating_failed", fmt.Sprintf("Cannot load actions: %v", err))
		return
	}

	cr.actions = append(cr.actions, acts...)
	for _, act := range acts {
		cr.ActionNames = append(cr.ActionNames, act.Name)
	}

	// no action, done!
	if len(cr.ActionNames) == 0 {
		cr.NewStatus("done", "no action")
		err = cr.saveAndFinish()
		return
	}

	cr.NewStatus("created", "")
	err = cr.Save()
	if err != nil {
		cr.NewStatus("saving_failed", fmt.Sprintf("Cannot save new run: %v", err))
		return
	}

	cr.NewStatus("queuing", "")
	err = cr.Trigger()
	if err != nil {
		cr.NewStatus("queuing_failed", fmt.Sprintf("Cannot queue actions: %v", err))
		err = cr.saveAndFinish()
		return
	}

	// succeed to trigger, change to queued
	cr.NewStatus("queued", "")
	err = cr.saveAndFinish()
	return
}

func (cr *CaseRun) saveAndFinish() (err error) {
	err = cr.Save()
	if err != nil {
		cr.NewStatus("saving_failed", fmt.Sprintf("Cannot save case run: %v", err))
	}
	return
}

func (cr *CaseRun) NewStatus(name, msg string) {
	cr.Stati = append(cr.Stati, Status{Name: name, At: time.Now(), Msg: msg})
}

func FindCaseRuns(path, name string) (runs []CaseRun, err error) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		return
	}
	defer session.Close()

	col := session.DB("legolas").C("runs")
	err = col.Find(bson.M{"path": path, "name": name}).All(&runs)
	return
}

// save CaseRun info to DB (not Redis)
func (cr *CaseRun) Save() (err error) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		return
	}
	defer session.Close()

	col := session.DB("legolas").C("runs")
	_, err = col.Upsert(bson.M{"uid": cr.Id}, *cr)
	return
}

// save actions into job queue
func (cr *CaseRun) Trigger() (err error) {
	rc, err := redisPool.Get()
	if err != nil {
		return
	}

	jobProxy := &CM.JobProxy{Queue: queueName, Rc: rc}

	var preActionName = ""
	for _, act := range cr.actions {
		err = jobProxy.Append(&CM.Job{
			CaseRunID: cr.Id,
			Name:      act.Name,
			Snippet:   act.Snippet,
			PreAction: preActionName,
		})
		if err != nil {
			L.Printf("inserting job failed: %v\n", err)
			break
		}
		preActionName = act.Name
	}
	return
}
