package run

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

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

func GetAll(tc *C.TestCase) (result []Run, err error) {
	err = col.Find(bson.M{"case_id": tc.Id}).All(&result)
	return
}

func New(caseId bson.ObjectId) *Run {
	return &Run{
		Id:     bson.NewObjectId(),
		CaseId: caseId,
	}
}

func (r *Run) Save() (err error) {
	_, err = col.Upsert(bson.M{"_id": r.Id}, *r)
	return
}

func (r *Run) Delete() (err error) {
	return col.Remove(bson.M{"_id": r.Id})
}
