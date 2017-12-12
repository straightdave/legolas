/*
Model of Action
*/
package action

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"

	T "legolas/common/models/template"
)

type Action struct {
	Id         bson.ObjectId          `json:"_id" bson:"_id"`
	CaseId     bson.ObjectId          `json:"case_id" bson:"case_id"`
	TemplateId bson.ObjectId          `json:"template_id" bson:"template_id"`
	Name       string                 `json:"name" bson:"name"`
	Desc       string                 `json:"desc" bson:"desc"`
	SeqNo      int                    `json:"seq_no" bson:"seq_no"`
	Params     map[string]interface{} `json:"params" bson:"params"`

	// for data mocking
	IsMocking bool                   `json:"is_mocking" bson:"is_mocking"`
	MockData  map[string]interface{} `json:"mock_data" bson:"mock_data"`

	Disabled  bool      `json:"disabled" bson:"disabled"`
	Removed   bool      `json:"removed" bson:"removed"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

func (a *Action) Json() ([]byte, error) {
	return bson.MarshalJSON(*a)
}

func (a *Action) JsonPretty() ([]byte, error) {
	return json.MarshalIndent(*a, "", "    ")
}

func FromJson(content []byte) (act Action, err error) {
	err = json.Unmarshal(content, &act)
	return
}

func (a *Action) ApplyTemplate(tpl *T.Template) {
	a.TemplateId = tpl.Id

	// set params with default value if it has,
	// or the nil value of user-defined types
	for k, v := range tpl.Params {
		if _, ok := a.Params[k]; ok {
			// if action has already has the param with that key,
			// ignore template's default value
			// TODO: need a type check?
			continue
		}

		if defaultValue, ok := v["default"]; ok {
			a.Params[k] = defaultValue
		} else {
			tRaw, ok := v["type"]
			if !ok {
				continue
			}

			if tStr, ok := tRaw.(string); ok {
				a.Params[k] = getDefault(tStr)
			}
		}
	}
}

func (a *Action) Snippet() (snippet string, err error) {
	T.SetCol(getMongo())
	tpl, err := T.GetOneById(a.TemplateId)
	if err != nil {
		return
	}
	snippet = tpl.Snippet
	return
}

func getDefault(typeStr string) interface{} {
	switch strings.ToLower(typeStr) {
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
