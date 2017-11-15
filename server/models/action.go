package models

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Action struct {
	CaseName string `bson:"case_name" json:"case_name"`
	CasePath string `bson:"case_path" json:"case_path"`
	Name     string `bson:"name" json:"name"`
	Desc     string `bson:"desc" json:"desc"`
	Snippet  string `bson:"snippt" json:"snippet"`
}

func (a *Action) FullName() string {
	return fmt.Sprintf("%s/%s#%s", a.CasePath, a.CaseName, a.Name)
}

func NewAction(cpath, cname, name string) *Action {
	return &Action{CasePath: cpath, CaseName: cname, Name: name}
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

func (a *Action) Save() error {
	session, err := mgo.Dial("localhost")
	if err != nil {
		return err
	}
	defer session.Close()

	col := session.DB("legolas").C("actions")
	return col.Insert(*a)
}

func UpdateActionOwner(cname, cpath, newName, newPath string) (err error) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		return
	}
	defer session.Close()

	var actions []Action
	col := session.DB("legolas").C("actions")
	err = col.Find(bson.M{"case_name": cname, "case_path": cpath}).All(&actions)
	if err != nil {
		return
	}

	// TODO: in a batch or transaction?
	for _, action := range actions {
		action.CaseName, action.CasePath = newName, newPath
		err = action.Save()
		if err != nil {
			break
		}
	}
	return
}
