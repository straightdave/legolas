package testcase

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	S "legolas/common/storage"
)

var (
	col *mgo.Collection
)

func SetCol(m *S.Mongo) {
	col = m.For("cases")
}

func getMongo() *S.Mongo {
	return &S.Mongo{
		Session: col.Database.Session,
	}
}

func GetAll(path string) (result []TestCase, err error) {
	err = col.Find(bson.M{"path": path}).All(&result)
	return
}

func GetOne(path, name string) (result TestCase, err error) {
	err = col.Find(bson.M{"path": path, "name": name}).One(&result)
	return
}

func GetOneById(id bson.ObjectId) (result TestCase, err error) {
	err = col.Find(bson.M{"_id": id}).One(&result)
	return
}

func New() (tc *TestCase) {
	return &TestCase{
		Id:     bson.NewObjectId(),
		Params: make(map[string]interface{}),
	}
}

func (tc *TestCase) Save() (err error) {
	if !tc.Id.Valid() {
		tc.Id = bson.NewObjectId()
	}
	_, err = col.Upsert(bson.M{"_id": tc.Id}, *tc)
	return
}

func (tc *TestCase) Delete() (err error) {
	return col.Remove(bson.M{"_id": tc.Id})
}
