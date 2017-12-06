package testcase

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"

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

// for server to list cases
// will get output of latest run for each case
func GetAllInTimeOrder(limit int) (result []TestCase, err error) {
	if limit < 1 {
		limit = 25
	}

	err = col.Find(bson.M{"removed": false}).Sort("-created_at").Limit(limit).All(&result)
	return
}

func GetAll(path string) (result []TestCase, err error) {
	err = col.Find(bson.M{"path": path, "removed": false}).All(&result)
	return
}

func GetFiltered(word string) (result []TestCase, err error) {
	err = col.Find(bson.M{"name": &bson.RegEx{Pattern: word, Options: "i"}, "removed": false}).Sort("-created_at").All(&result)
	return
}

func GetOne(path, name string) (result TestCase, err error) {
	err = col.Find(bson.M{"path": path, "name": name, "removed": false}).One(&result)
	return
}

func GetOneById(id bson.ObjectId) (result TestCase, err error) {
	// even the deleted one
	err = col.Find(bson.M{"_id": id}).One(&result)
	return
}

func New() (tc *TestCase) {
	return &TestCase{
		Id:        bson.NewObjectId(),
		Params:    make(map[string]interface{}),
		CreatedAt: time.Now(),
	}
}

func (tc *TestCase) Save() (err error) {
	if !tc.Id.Valid() {
		tc.Id = bson.NewObjectId()
	}
	if tc.CreatedAt == (time.Time{}) {
		tc.CreatedAt = time.Now()
	}
	tc.UpdatedAt = time.Now()
	_, err = col.Upsert(bson.M{"_id": tc.Id}, *tc)
	return
}

func (tc *TestCase) Delete() (err error) {
	if !tc.Id.Valid() {
		err = errors.New("Test case Id is invalid")
		return
	}

	tc.Removed = true
	_, err = col.Upsert(bson.M{"_id": tc.Id}, *tc)
	return
}
