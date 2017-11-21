package testcase

import (
	"fmt"
	"testing"
)

func TestCreateTestCase(t *testing.T) {
	// create basic and save
	c := NewTestCase("$case/path", "case-1", "this is a case")
	if err := c.Save(); err != nil {
		t.Fatalf("cannot save new case: %v\n", err)
	}

	datap, err := c.JsonPretty()
	if err != nil {
		t.Fatalf("converting to json failed: %v\n", err)
	}
	fmt.Println(string(datap))

	// add params
	c.AddParam("name", "dave")
	c.AddParam("age", 18)

	if err := c.Save(); err != nil {
		t.Fatalf("cannot update & save case: %v\n", err)
	}

	datap, err = c.JsonPretty()
	if err != nil {
		t.Fatalf("converting to json failed: %v\n", err)
	}
	fmt.Println(string(datap))

	// delete
	if err := c.Delete(); err != nil {
		t.Fatalf("cannot delete template: %v\n", err)
	}
}
