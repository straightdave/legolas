/*
Model of Action
*/
package action

import (
	"encoding/json"
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"legolas/common/config"
	"legolas/common/models/template"
)

type Action struct {
	CasePath     string                 `json:"case_path" bson:"case_path"`
	CaseName     string                 `json:"case_name" bson:"case_name"`
	Name         string                 `json:"name" bson:"name"`
	TemplatePath string                 `json:"template_path" bson:"template_path"`
	TemplateName string                 `json:"template_name" bson:"template_name"`
	Snippet      string                 `json:"snippet" bson:"snippet"`
	Params       map[string]interface{} `json:"params" bson:"params"`
}

func (a *Action) Json() ([]byte, error) {
	return json.Marshal(*a)
}

func (a *Action) JsonPretty() ([]byte, error) {
	return json.MarshalIndent(*a, "", "    ")
}

func (a *Action) isValid() bool {
	return a.CasePath != "" && a.CaseName != "" && a.Name != ""
}

func getDefaultValueOfType(t string) interface{} {
	switch t {
	case "text", "string":
		return ""
	case "float", "decimal", "double":
		return 0.0
	case "number", "int":
		return 0
	case "bool", "boolean":
		return false
	default:
		return nil
	}
}

func NewAction(cpath, cname, name string) Action {
	return Action{
		CasePath: cpath,
		CaseName: cname,
		Name:     name,
		Params:   make(map[string]interface{}),
	}
}

func (a *Action) ApplyTemplate(tpath, tname string) (err error) {
	a.TemplatePath = tpath
	a.TemplateName = tname

	tpl, err := template.GetTemplate(tpath, tname)
	if err != nil {
		return
	}

	a.Snippet = tpl.Snippet

	// set params with default value if it has,
	// or the nil value of user-defined types
	for k, v := range tpl.Params {
		if defaultValue, ok := v["default"]; ok {
			a.Params[k] = defaultValue
		} else {
			a.Params[k] = getDefaultValueOfType(v["type"].(string))
		}
	}
	return
}

func (a *Action) SetParam(pname string, pvalue interface{}) {
	a.Params[pname] = pvalue
}

func GetActions(cpath, cname string) (result []Action, err error) {
	session, err := mgo.Dial(config.MongoHost)
	if err != nil {
		return
	}
	defer session.Close()

	err = session.DB("legolas").C("actions").Find(bson.M{"case_path": cpath, "case_name": cname}).All(&result)
	return
}

func (a *Action) Save() (err error) {
	if !a.isValid() {
		return errors.New("action case_path or case_name or name is blank")
	}

	session, err := mgo.Dial(config.MongoHost)
	if err != nil {
		return
	}
	defer session.Close()

	_, err = session.DB("legolas").C("actions").Upsert(bson.M{"case_path": a.CasePath, "case_name": a.CaseName, "name": a.Name}, *a)
	return
}

func (a *Action) Delete() (err error) {
	if !a.isValid() {
		return errors.New("action case_path or case_name or name is blank")
	}

	session, err := mgo.Dial(config.MongoHost)
	if err != nil {
		return
	}
	defer session.Close()

	err = session.DB("legolas").C("actions").Remove(bson.M{"case_path": a.CasePath, "case_name": a.CaseName, "name": a.Name})
	return
}
