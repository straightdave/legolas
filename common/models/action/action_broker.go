package action

import (
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
	err = col.Find(bson.M{"case_path": cpath, "case_name": cname}).All(&result)
	return
}

func GetAllByCaseId(caseId bson.ObjectId) (result []Action, err error) {
	err = col.Find(bson.M{"case_id": caseId}).Sort("seq_no").All(&result)
	return
}

func GetOne(cpath, cname, name string) (result Action, err error) {
	err = col.Find(bson.M{"case_path": cpath, "case_name": cname, "name": name}).One(&result)
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
	_, err = col.Upsert(bson.M{"_id": a.Id}, *a)
	return
}

func (a *Action) Delete() (err error) {
	err = col.Remove(bson.M{"_id": a.Id})
	return
}
