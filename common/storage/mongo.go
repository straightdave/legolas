package storage

import (
	"gopkg.in/mgo.v2"
	L "log"

	"legolas/common/config"
)

type Mongo struct {
	Session *mgo.Session
}

var (
	ses *mgo.Session
)

func init() {
	var err error
	ses, err = mgo.Dial(config.MongoHost)
	if err != nil {
		L.Fatalf("failed to init mongo: %v\n", err)
	}
}

func AskForMongo() (m *Mongo) {
	return &Mongo{
		Session: ses.Copy(),
	}
}

func (m *Mongo) Close() {
	m.Session.Close()
}

func (m *Mongo) For(col string) *mgo.Collection {
	return m.Session.DB("legolas").C(col)
}
