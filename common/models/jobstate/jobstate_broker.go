package jobstate

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	A "legolas/common/models/action"
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
	err = col.Find(bson.M{"run_id": runId}).Sort("created_at").All(&result)
	return
}

func GetAllByRunIdStr(runId string) (result []JobState, err error) {
	rid := bson.ObjectIdHex(runId)
	if !rid.Valid() {
		err = errors.New("Invalid run id")
		return
	}
	err = col.Find(bson.M{"run_id": rid}).Sort("created_at").All(&result)
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
	if js.ActionName == "" {
		A.SetCol(getMongo())
		act, err := A.GetOneById(js.ActionId)
		if err != nil {
			return err
		}
		js.ActionName = act.Name
	}
	_, err = col.Upsert(bson.M{"run_id": js.RunId, "action_id": js.ActionId}, *js)
	return
}
