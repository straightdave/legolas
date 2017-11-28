package template

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	json := `
    {
        "_id": "5a1c1e092750834938a36189",
        "path": "$template/path/path",
        "name": "my-template",
        "desc": "this is a test template",
        "snippet": "print(\"hello world\")",
        "params": {
            "name": {
                "type": "text",
                "required": true
            },
            "age": {
                "type": "number",
                "required": true
            }
        }
    }`

	tpl, err := FromJson([]byte(json))
	if err != nil {
		t.Errorf("cannot parse json: %v\n", err)
	}

	if !tpl.Id.Valid() {
		t.Errorf("tpl.Id is not valid\n")
	}

	if tpl.Id != bson.ObjectIdHex("5a1c1e092750834938a36189") {
		t.Errorf("tpl.id is not objectId\n")
	}

	c, _ := tpl.JsonPretty()
	fmt.Println(string(c))
}

func TestMarshal(t *testing.T) {
	tpl := New()
	c, _ := tpl.Json()
	fmt.Println(string(c))
}

func TestConvertTemplateToJson(t *testing.T) {
	tpl := &Template{
		Id:      bson.NewObjectId(),
		Path:    "$template/path",
		Name:    "template-1",
		Desc:    "this is a template. wahaha.",
		Snippet: `print("hello world")`,
		Params: map[string]map[string]interface{}{
			"param-1": {
				"type":     "string",
				"required": true,
			},
			"param-2": {
				"type":     "number",
				"required": false,
			},
		},
	}

	data, err := tpl.Json()
	if err != nil {
		t.Fatalf("converting to json failed: %v\n", err)
	}
	fmt.Println(string(data))

	datap, err := tpl.JsonPretty()
	if err != nil {
		t.Fatalf("converting to json failed: %v\n", err)
	}
	fmt.Println(string(datap))

	tpl1, _ := FromJson(data)
	tpl2, _ := FromJson(datap)
	if tpl.Id != tpl1.Id && tpl1.Id != tpl2.Id {
		fmt.Println("converting json and jsonPretty back is failed")
		t.Fail()
	}
	fmt.Println("they are equal")
}
