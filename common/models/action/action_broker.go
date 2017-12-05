package action

import (
	"errors"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	TC "legolas/common/models/testcase"
	S "legolas/common/storage"
)

var (
	col *mgo.Collection
)

func SetCol(m *S.Mongo) {
	col = m.For("actions")
}

// internal use
func getMongo() *S.Mongo {
	return &S.Mongo{
		Session: col.Database.Session,
	}
}

func GetAll(cpath, cname string) (result []Action, err error) {
	TC.SetCol(getMongo())
	tc, err := TC.GetOne(cpath, cname)
	if err != nil {
		return
	}
	err = col.Find(bson.M{"case_id": tc.Id, "removed": false}).All(&result)
	return
}

func GetAllByCaseId(caseId bson.ObjectId) (result []Action, err error) {
	err = col.Find(bson.M{"case_id": caseId, "removed": false}).Sort("seq_no").All(&result)
	return
}

func GetAllByCaseIdStr(caseId string) (result []Action, err error) {
	_cid := bson.ObjectIdHex(caseId)
	err = col.Find(bson.M{"case_id": _cid, "removed": false}).Sort("seq_no").All(&result)
	return
}

func GetOne(cpath, cname, name string) (result Action, err error) {
	TC.SetCol(getMongo())
	tc, err := TC.GetOne(cpath, cname)
	if err != nil {
		return
	}
	err = col.Find(bson.M{"case_id": tc.Id, "name": name, "removed": false}).One(&result)
	return
}

func GetOneById(actionId bson.ObjectId) (result Action, err error) {
	err = col.Find(bson.M{"_id": actionId}).One(&result)
	return
}

func CountByCase(caseId bson.ObjectId) (result int, err error) {
	return col.Find(bson.M{"case_id": caseId, "removed": false}).Count()
}

func MaxSeqNoInCase(caseId bson.ObjectId) (result int, err error) {
	result = 0
	var act Action
	err = col.Find(bson.M{"case_id": caseId, "removed": false}).Sort("-seq_no").One(&act)
	if err != nil {
		err = errors.New(fmt.Sprintf("case_id: %v, sortby seq_no failed: %v", caseId, err))
	}
	result = act.SeqNo
	return
}

func New() *Action {
	return &Action{
		Id:     bson.NewObjectId(),
		Params: make(map[string]interface{}),
	}
}

func (a *Action) Save() (err error) {
	if a.CaseId.Valid() && a.TemplateId.Valid() {
		if !a.Id.Valid() {
			a.Id = bson.NewObjectId()
		}

		if a.SeqNo < 1 {
			maxSeqNo, _ := MaxSeqNoInCase(a.CaseId) //TODO: investigate: why blank case (no action) will fail to sort
			a.SeqNo = maxSeqNo + 5
		}
		_, err = col.Upsert(bson.M{"_id": a.Id}, *a)
	} else {
		err = errors.New("Case Id or Template Id is invalid")
	}
	return
}

func (a *Action) Delete() (err error) {
	if !a.Id.Valid() {
		err = errors.New("Action Id is invalid")
		return
	}
	a.Removed = true
	return a.Save()
}
