package action

import (
	"fmt"
	"legolas/common/models/template"
	"testing"
)

func TestConvertActionToJson(t *testing.T) {
	a := &Action{
		CasePath:     "$case/path",
		CaseName:     "case-1",
		Name:         "action-1",
		TemplatePath: "$template/path",
		TemplateName: "template-1",
		Params: map[string]interface{}{
			"param-1": "hello",
			"param-2": 123,
		},
	}

	data, err := a.Json()
	if err != nil {
		t.Fatalf("converting to json failed: %v\n", err)
	}
	fmt.Println(string(data))

	datap, err := a.JsonPretty()
	if err != nil {
		t.Fatalf("converting to json (pretty) failed: %v\n", err)
	}
	fmt.Println(string(datap))
}

func TestAddNewAction(t *testing.T) {
	// create a template
	tpl := template.NewTemplate("$template/path", "template-1", "this is a template")
	tpl.Snippet = `print("hello world")`
	tpl.AddParam("name", "text", "dave", true)
	tpl.AddParam("age", "number", 18, true)
	if err := tpl.Save(); err != nil {
		t.Fatalf("failed to create new template: %v\n", err)
	}

	// create action of template
	a := NewAction("$case/path", "case-1", "action-1")
	a.ApplyTemplate("$template/path", "template-1")

	// save new created action
	if err := a.Save(); err != nil {
		t.Fatalf("cannot save new action:%v\n", err)
	}

	datap, err := a.JsonPretty()
	if err != nil {
		t.Fatalf("converting to json (pretty) failed: %v\n", err)
	}
	fmt.Println(string(datap))

	// delete action and template
	if err := a.Delete(); err != nil {
		t.Fatalf("cannot delete action: %v\n", err)
	}

	if err := tpl.Delete(); err != nil {
		t.Fatalf("cannot delete template: %v\n", err)
	}
}
