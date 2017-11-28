package template

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
)

type Template struct {
	Id      bson.ObjectId                     `json:"_id" bson:"_id"`
	Path    string                            `json:"path" bson:"path"`
	Name    string                            `json:"name" bson:"name"`
	Desc    string                            `json:"desc" bson:"desc"`
	Params  map[string]map[string]interface{} `json:"params" bson:"params"`
	Snippet string                            `json:"snippet" bson:"snippet"`
}

func (tpl *Template) Json() ([]byte, error) {
	return bson.MarshalJSON(*tpl)
}

func (tpl *Template) JsonPretty() ([]byte, error) {
	return json.MarshalIndent(*tpl, "", "    ")
}

func FromJson(content []byte) (tpl Template, err error) {
	err = bson.UnmarshalJSON(content, &tpl)
	return
}

func (tpl *Template) AddParam(pname, ptype string, pdefault interface{}, prequired bool) {
	if tpl.Params == nil {
		tpl.Params = make(map[string]map[string]interface{})
	}

	tpl.Params[pname] = map[string]interface{}{
		"type":     ptype,
		"required": prequired,
		"default":  pdefault,
	}
}
