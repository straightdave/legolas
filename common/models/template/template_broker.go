package template

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
	col = m.For("templates")
}

// used internally to fetch mongo session
func getMongo() *S.Mongo {
	return &S.Mongo{
		Session: col.Database.Session,
	}
}

func GetAll(path string) (result []Template, err error) {
	err = col.Find(bson.M{"path": path, "removed": false}).All(&result)
	return
}

func GetOne(path, name string) (result Template, err error) {
	err = col.Find(bson.M{"path": path, "name": name, "removed": false}).One(&result)
	return
}

func GetOneById(id bson.ObjectId) (result Template, err error) {
	err = col.Find(bson.M{"_id": id}).One(&result)
	return
}

func New() (tpl *Template) {
	return &Template{
		Id:     bson.NewObjectId(),
		Params: make(map[string]map[string]interface{}),
	}
}

func (tpl *Template) Save() (err error) {
	if !tpl.Id.Valid() {
		tpl.Id = bson.NewObjectId()
	}
	_, err = col.Upsert(bson.M{"_id": tpl.Id}, *tpl)
	return
}

func (tpl *Template) Delete() (err error) {
	if !tpl.Id.Valid() {
		err = errors.New("Template Id is invalid")
		return
	}

	tpl.Removed = true
	_, err = col.Upsert(bson.M{"_id": tpl.Id}, *tpl)
	return
}
