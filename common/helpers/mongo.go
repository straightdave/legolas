package helpers

import (
	"gopkg.in/mgo.v2"
	L "log"

	"legolas/common/config"
)

type Mongo struct {
	session *mgo.Session
}

func NewMongo() *Mongo {
	sess, err := mgo.Dial(config.MongoHost)
	if err != nil {
		L.Fatalln("failed to create new Mongo")
	}

	return &Mongo{
		session: sess,
	}
}

func (m *Mongo) Close() {
	m.session.Close()
}

func (m *Mongo) Get() *mgo.Session {
	return m.session.Copy()
}
