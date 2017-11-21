package template

import (
	"fmt"
	"testing"
)

func TestConvertTemplateToJson(t *testing.T) {
	tpl := &Template{
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
}

func TestGetTemplate(t *testing.T) {
	// create basic and save
	tpl := NewTemplate("$temp/path", "template-1", "this is a template")
	if err := tpl.Save(); err != nil {
		t.Fatalf("cannot save new template: %v\n", err)
	}

	datap, err := tpl.JsonPretty()
	if err != nil {
		t.Fatalf("converting to json failed: %v\n", err)
	}
	fmt.Println(string(datap))

	// add more fields and save
	tpl.Snippet = `print("hello world")`
	tpl.Params = map[string]map[string]interface{}{
		"p1": {
			"type":     "text",
			"required": true,
			"default":  "hello",
		},
		"p2": {
			"type":     "number",
			"required": false,
			"default":  123,
		},
	}
	if err := tpl.Save(); err != nil {
		t.Fatalf("cannot update & save template: %v\n", err)
	}

	datap, err = tpl.JsonPretty()
	if err != nil {
		t.Fatalf("converting to json failed: %v\n", err)
	}
	fmt.Println(string(datap))

	// delete
	if err := tpl.Delete(); err != nil {
		t.Fatalf("cannot delete template: %v\n", err)
	}
}

func TestAddParam(t *testing.T) {
	// create basic and save
	tpl := NewTemplate("$temp/path", "template-1", "this is a template")
	if err := tpl.Save(); err != nil {
		t.Fatalf("cannot save new template: %v\n", err)
	}

	datap, err := tpl.JsonPretty()
	if err != nil {
		t.Fatalf("converting to json failed: %v\n", err)
	}
	fmt.Println(string(datap))

	// add params
	tpl.AddParam("name", "text", "dave", true)
	tpl.AddParam("age", "number", 18, false)

	// save
	if err := tpl.Save(); err != nil {
		t.Fatalf("cannot update & save template: %v\n", err)
	}

	datap, err = tpl.JsonPretty()
	if err != nil {
		t.Fatalf("converting to json failed: %v\n", err)
	}
	fmt.Println(string(datap))

	// delete
	if err := tpl.Delete(); err != nil {
		t.Fatalf("cannot delete template: %v\n", err)
	}
}
