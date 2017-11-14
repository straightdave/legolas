package models

import (
	"fmt"
	"testing"
)

func TestSaveNew(t *testing.T) {
	c := &Case{
		Path: "group-a",
		Name: "case-1",
		Desc: "this is a test case no.1",
	}

	if c.FullName() != "group-a/case-1" {
		t.Errorf("fullname %s not expected\n", c.FullName())
	}

	err := DeleteCase("group-a", "case-1")
	if err != nil {
		t.Errorf("cannot delete: %v\n", err)
	}

	err = c.Save()
	if err != nil {
		t.Errorf("cannot save: %v\n", err)
	}
}

func TestGet(t *testing.T) {
	cases, err := FindAllCases()
	if err != nil {
		t.Errorf("cannot get all: %v", err)
	} else {
		fmt.Printf("get %d cases\n", len(cases))
	}

	err = DeleteCase("path1", "case1")
	if err != nil {
		t.Errorf("cannot delete: %v\n", err)
	}

	err = NewCase("path1", "case1", "test").Save()
	if err != nil {
		t.Errorf("cannot create new one: %v\n", err)
	}

	c, err := FindCase("path1", "case1")
	if err != nil {
		t.Errorf("can not get one: %v\n", err)
	}

	if c.Name != "case1" || c.Path != "path1" {
		t.Errorf("get wrong one\n")
	}
}
