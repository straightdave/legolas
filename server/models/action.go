package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Action struct {
	CaseName string `bson:"case_name" json:"case_name"`
	CasePath string `bson:"case_path" json:"case_path"`
	Name     string `bson:"name" json:"name"`
}

func FindActions(cname, cpath string) (result []Action, err error) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		return nil, err
	}
	defer session.Close()

	col := session.DB("legolas").C("actions")
	err = col.Find(bson.M{"case_name": cname, "case_path": cpath}).All(&result)
	return
}

func FindAllActions() (result []Action, err error) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		return nil, err
	}
	defer session.Close()

	col := session.DB("legolas").C("actions")
	err = col.Find(nil).All(&result)
	return
}
