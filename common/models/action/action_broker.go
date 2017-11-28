package action

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

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
	err = col.Find(bson.M{"case_path": cpath, "case_name": cname, "removed": false}).All(&result)
	return
}

func GetAllByCaseId(caseId bson.ObjectId) (result []Action, err error) {
	err = col.Find(bson.M{"case_id": caseId, "removed": false}).Sort("seq_no").All(&result)
	return
}

func GetOne(cpath, cname, name string) (result Action, err error) {
	err = col.Find(bson.M{"case_path": cpath, "case_name": cname, "name": name, "removed": false}).One(&result)
	return
}

func GetOneById(actionId bson.ObjectId) (result Action, err error) {
	err = col.Find(bson.M{"_id": actionId}).One(&result)
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
	_, err = col.Upsert(bson.M{"_id": a.Id}, *a)
	return
}
