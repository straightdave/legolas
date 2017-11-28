package jobstate

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	S "legolas/common/storage"
)

var (
	col *mgo.Collection
)

func SetCol(m *S.Mongo) {
	col = m.For("jobstates")
}

func getMongo() *S.Mongo {
	return &S.Mongo{
		Session: col.Database.Session,
	}
}

func GetAll(runId bson.ObjectId) (result []JobState, err error) {
	err = col.Find(bson.M{"run_id": runId}).All(&result)
	return
}

func GetOne(runId, actionId bson.ObjectId) (result JobState, err error) {
	err = col.Find(bson.M{"run_id": runId, "action_id": actionId}).One(&result)
	return
}

func New(runId, actionId bson.ObjectId) *JobState {
	return &JobState{
		RunId:    runId,
		ActionId: actionId,
	}
}

func (js *JobState) Save() (err error) {
	_, err = col.Upsert(bson.M{"run_id": js.RunId, "action_id": js.ActionId}, *js)
	return
}
