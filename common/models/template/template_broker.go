package template

import (
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
	err = col.Find(bson.M{"path": path}).All(&result)
	return
}

func GetOne(path, name string) (result Template, err error) {
	err = col.Find(bson.M{"path": path, "name": name}).One(&result)
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
	_, err = col.Upsert(bson.M{"_id": tpl.Id}, *tpl)
	return
}

func (tpl *Template) Delete() (err error) {
	return col.Remove(bson.M{"_id": tpl.Id})
}
