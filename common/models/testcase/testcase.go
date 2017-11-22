/*
Model of test case
*/
package testcase

import (
	"encoding/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"errors"
	"legolas/common/config"
)

type TestCase struct {
	Path   string                 `json:"path" bson:"path"`
	Name   string                 `json:"name" bson:"name"`
	Desc   string                 `json:"desc" bson:"desc"`
	Params map[string]interface{} `json:"params" bson:"params"`
}

func (c *TestCase) Json() ([]byte, error) {
	return bson.MarshalJSON(*c)
}

func (c *TestCase) JsonPretty() ([]byte, error) {
	return json.MarshalIndent(*c, "", "    ")
}

func (c *TestCase) isValid() bool {
	return c.Path != "" && c.Name != ""
}

func NewTestCase(path, name, desc string) TestCase {
	return TestCase{
		Path:   path,
		Name:   name,
		Desc:   desc,
		Params: make(map[string]interface{}),
	}
}

func FromJson(content []byte) (ca TestCase, err error) {
	ca.Params = make(map[string]interface{}) // in case not initialized
	err = json.Unmarshal(content, &ca)
	return
}

func (c *TestCase) AddParam(pname string, pvalue interface{}) {
	c.Params[pname] = pvalue
}

func GetTestCases(path string) (result []TestCase, err error) {
	session, err := mgo.Dial(config.MongoHost)
	if err != nil {
		return
	}
	defer session.Close()

	err = session.DB("legolas").C("cases").Find(bson.M{"path": path}).All(&result)
	return
}

func GetTestCase(path, name string) (result TestCase, err error) {
	session, err := mgo.Dial(config.MongoHost)
	if err != nil {
		return
	}
	defer session.Close()

	err = session.DB("legolas").C("cases").Find(bson.M{"path": path, "name": name}).One(&result)
	return
}

func (c *TestCase) Save() (err error) {
	if !c.isValid() {
		return errors.New("case path or name is blank")
	}

	session, err := mgo.Dial(config.MongoHost)
	if err != nil {
		return
	}
	defer session.Close()

	_, err = session.DB("legolas").C("cases").Upsert(bson.M{"path": c.Path, "name": c.Name}, *c)
	return
}

func (c *TestCase) Delete() (err error) {
	if !c.isValid() {
		return errors.New("case path or name is blank")
	}

	session, err := mgo.Dial(config.MongoHost)
	if err != nil {
		return
	}
	defer session.Close()

	err = session.DB("legolas").C("cases").Remove(bson.M{"path": c.Path, "name": c.Name})
	return
}
