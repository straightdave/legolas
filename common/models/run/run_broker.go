package run

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"

	A "legolas/common/models/action"
	J "legolas/common/models/job"
	C "legolas/common/models/testcase"
	S "legolas/common/storage"
)

var (
	col *mgo.Collection
)

func SetCol(m *S.Mongo) {
	col = m.For("runs")
}

// internal use
func getMongo() *S.Mongo {
	return &S.Mongo{
		Session: col.Database.Session,
	}
}

func GetOne(runId bson.ObjectId) (result Run, err error) {
	err = col.Find(bson.M{"_id": runId}).One(&result)
	return
}

func GetAll(tc *C.TestCase) (result []Run, err error) {
	err = col.Find(bson.M{"case_id": tc.Id}).All(&result)
	return
}

func GetAllByCaseId(caseId bson.ObjectId) (result []Run, err error) {
	err = col.Find(bson.M{"case_id": caseId}).Sort("-started_at").All(&result)
	return
}

func GetAllByCaseIdStr(caseId string) (result []Run, err error) {
	cid := bson.ObjectIdHex(caseId)
	if !cid.Valid() {
		err = errors.New("invalid case Id")
		return
	}
	err = col.Find(bson.M{"case_id": cid}).Sort("-started_at").All(&result)
	return
}

func New(caseId bson.ObjectId) *Run {
	return &Run{
		Id:     bson.NewObjectId(),
		CaseId: caseId,
	}
}

func (r *Run) Save() (err error) {
	if r.StartedAt == (time.Time{}) {
		r.StartedAt = time.Now()
	}
	_, err = col.Upsert(bson.M{"_id": r.Id}, *r)
	return
}

func (r *Run) Delete() (err error) {
	return col.Remove(bson.M{"_id": r.Id})
}

func InvokeByCaseIdStr(caseId string) (run *Run, err error) {
	cid := bson.ObjectIdHex(caseId)
	if !cid.Valid() {
		err = errors.New("invalid case Id")
		return
	}

	run = New(cid)

	A.SetCol(getMongo())
	actions, err := A.GetAllByCaseId(cid)
	if err != nil {
		return
	}

	// serialized pushing for orders
	var prev bson.ObjectId
	for _, act := range actions {
		job := &J.Job{
			RunId:        run.Id,
			ActionId:     act.Id,
			PrevActionId: prev,
		}
		err = J.Append(job)
		if err != nil {
			return
		}
		prev = act.Id
	}
	err = run.Save()
	return
}
