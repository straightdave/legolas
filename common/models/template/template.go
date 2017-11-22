package template

import (
	"encoding/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"errors"
	"legolas/common/config"
)

type Template struct {
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
	tpl.Params = make(map[string]map[string]interface{}) // in case not initialized
	err = json.Unmarshal(content, &tpl)
	if err != nil {
		return
	}

	if !tpl.isValid() {
		err = errors.New("template loaded is invalid.")
	}

	return
}

func GetTemplates(path string) (result []Template, err error) {
	session, err := mgo.Dial(config.MongoHost)
	if err != nil {
		return
	}
	defer session.Close()
	err = session.DB("legolas").C("templates").Find(bson.M{"path": path}).All(&result)
	return
}

func GetTemplate(path, name string) (result Template, err error) {
	session, err := mgo.Dial(config.MongoHost)
	if err != nil {
		return
	}
	defer session.Close()
	err = session.DB("legolas").C("templates").Find(bson.M{"path": path, "name": name}).One(&result)
	return
}

func (tpl *Template) isValid() bool {
	return tpl.Path != "" && tpl.Name != ""
}

func NewTemplate(path, name, desc string) Template {
	return Template{
		Path:    path,
		Name:    name,
		Desc:    desc,
		Params:  make(map[string]map[string]interface{}),
		Snippet: "",
	}
}

func (tpl *Template) AddParam(pname, ptype string, pdefault interface{}, prequired bool) {
	tpl.Params[pname] = map[string]interface{}{
		"type":     ptype,
		"required": prequired,
		"default":  pdefault,
	}
}

func (tpl *Template) Save() (err error) {
	if !tpl.isValid() {
		return errors.New("template path or name is blank")
	}

	session, err := mgo.Dial(config.MongoHost)
	if err != nil {
		return
	}
	defer session.Close()

	_, err = session.DB("legolas").C("templates").Upsert(bson.M{"path": tpl.Path, "name": tpl.Name}, *tpl)
	return
}

func (tpl *Template) Delete() (err error) {
	if !tpl.isValid() {
		return errors.New("template path or name is blank")
	}

	session, err := mgo.Dial(config.MongoHost)
	if err != nil {
		return
	}
	defer session.Close()

	err = session.DB("legolas").C("templates").Remove(bson.M{"path": tpl.Path, "name": tpl.Name})
	return
}
