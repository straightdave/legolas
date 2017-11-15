package models

import (
	"errors"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Case struct {
	Path string `json:"path"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}

func (c *Case) FullName() string {
	return fmt.Sprintf("%s/%s", c.Path, c.Name)
}

func NewCase(path, name, desc string) *Case {
	return &Case{Path: path, Name: name, Desc: desc}
}

func FindAllCases() (result []Case, err error) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		return nil, err
	}
	defer session.Close()

	col := session.DB("legolas").C("cases")
	err = col.Find(bson.M{}).All(&result)
	return
}

func FindCase(path, name string) (result *Case, err error) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		return nil, err
	}
	defer session.Close()

	col := session.DB("legolas").C("cases")
	var res Case
	err = col.Find(bson.M{"path": path, "name": name}).One(&res)
	result = &res
	return
}

func (c *Case) UpdateTo(newCase *Case) (err error) {
	session, err := mgo.Dial("localhost")
	if err != nil {
		return
	}
	defer session.Close()

	// update all actions if necessary
	if c.FullName() != newCase.FullName() {
		err = UpdateActionOwner(c.Name, c.Path, newCase.Name, newCase.Path)
		if err != nil {
			return
		}
	}

	col := session.DB("legolas").C("cases")

	// delete old one
	err = col.Remove(bson.M{"path": c.Path, "name": c.Name})
	if err != nil {
		return
	}

	// save new one
	return col.Insert(*newCase)
}

func (c *Case) Save() error {
	session, err := mgo.Dial("localhost")
	if err != nil {
		return err
	}
	defer session.Close()

	col := session.DB("legolas").C("cases")
	n, err := col.Find(bson.M{"path": c.Path, "name": c.Name}).Count()
	if err != nil {
		return err
	}

	if n > 0 {
		return errors.New("duplicated full name: " + c.FullName())
	}

	return col.Insert(*c)
}

func DeleteCase(path, name string) error {
	session, err := mgo.Dial("localhost")
	if err != nil {
		return err
	}
	defer session.Close()

	// TODO: need to remove all actions as well

	col := session.DB("legolas").C("cases")
	return col.Remove(bson.M{"path": path, "name": name})
}
