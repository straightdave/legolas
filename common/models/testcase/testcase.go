/*
Model of test case
*/
package testcase

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type TestCase struct {
	Id        bson.ObjectId          `json:"_id" bson:"_id"`
	Path      string                 `json:"path" bson:"path"`
	Name      string                 `json:"name" bson:"name"`
	Desc      string                 `json:"desc" bson:"desc"`
	Params    map[string]interface{} `json:"params" bson:"params"`
	Disabled  bool                   `json:"disabled" bson:"disabled"`
	Removed   bool                   `json:"removed" bson:"removed"`
	CreatedAt time.Time              `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time              `json:"updated_at" bson:"updated_at"`
}

func (tc *TestCase) Json() ([]byte, error) {
	return bson.MarshalJSON(*tc)
}

func (tc *TestCase) JsonPretty() ([]byte, error) {
	return json.MarshalIndent(*tc, "", "    ")
}

func FromJson(content []byte) (tc TestCase, err error) {
	// in this case, the Id is blank
	// so in Save()/Delete() we should check about Id
	err = json.Unmarshal(content, &tc)
	return
}

func (tc *TestCase) AddParam(pname string, pvalue interface{}) {
	tc.Params[pname] = pvalue
}
