/*
Model of test case
*/
package testcase

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
)

type TestCase struct {
	Id     bson.ObjectId          `json:"_id" bson:"_id"`
	Path   string                 `json:"path" bson:"path"`
	Name   string                 `json:"name" bson:"name"`
	Desc   string                 `json:"desc" bson:"desc"`
	Params map[string]interface{} `json:"params" bson:"params"`
}

func (tc *TestCase) Json() ([]byte, error) {
	return bson.MarshalJSON(*tc)
}

func (tc *TestCase) JsonPretty() ([]byte, error) {
	return json.MarshalIndent(*tc, "", "    ")
}

func FromJson(content []byte) (tc TestCase, err error) {
	err = json.Unmarshal(content, &tc)
	return
}

func (tc *TestCase) AddParam(pname string, pvalue interface{}) {
	tc.Params[pname] = pvalue
}
