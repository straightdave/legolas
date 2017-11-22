package template

import (
	"fmt"
	"testing"
)

func TestTemplateFromJson(t *testing.T) {
	json := `
    {
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
		t.Fatalf("cannot parse json: %v\n", err)
	}

	c, _ := tpl.JsonPretty()
	fmt.Println(string(c))
}
