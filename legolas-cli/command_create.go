package main

import (
	"fmt"
	"io/ioutil"

	A "legolas/common/models/action"
	T "legolas/common/models/template"
	TC "legolas/common/models/testcase"
)

func createCommand(args []string) {
	if len(args) < 3 {
		fmt.Println("...what? (case, action, template)")
		return
	}

	switch args[2] {
	case "template":
		if len(args) < 5 {
			fmt.Println("lacks of arguments: need def file name and snippet file name")
			return
		}
		createTemplate(args[3], args[4])

	case "action":
		if len(args) < 4 {
			fmt.Println("lacks of arguments: need def file name ")
			return
		}
		createAction(args[3])

	case "case":
		if len(args) < 4 {
			fmt.Println("lacks of arguments: need def file name ")
			return
		}
		createCase(args[3])

	default:
		fmt.Println("unknown object type " + args[2])
	}
}

func createTemplate(defJsonFile, snippetFile string) {
	defjson, err := ioutil.ReadFile(defJsonFile)
	if err != nil {
		fmt.Printf("failed to read definition file of template: %v\n", err)
		return
	}

	tpl, err := T.FromJson(defjson)
	if err != nil {
		fmt.Printf("failed to parse definition of template: %v\n", err)
		return
	}

	snippet, err := ioutil.ReadFile(snippetFile)
	if err != nil {
		fmt.Printf("failed to read snippet file of template: %v\n", err)
		return
	}

	tpl.Snippet = string(snippet)
	if err := tpl.Save(); err != nil {
		fmt.Printf("failed to create new template: %v\n", err)
		return
	}
	fmt.Println("succeed creating template:")
	o, _ := tpl.JsonPretty()
	fmt.Println(string(o))
}

func createAction(defJsonFile string) {
	defjson, err := ioutil.ReadFile(defJsonFile)
	if err != nil {
		fmt.Printf("failed to read definition file of action: %v\n", err)
		return
	}

	act, err := A.FromJson(defjson)
	if err != nil {
		fmt.Printf("failed to parse definition of action: %v\n", err)
		return
	}

	tpl, err := T.GetOneById(act.TemplateId)
	if err != nil {
		fmt.Printf("failed to get template data: %v\n", err)
		return
	}

	act.ApplyTemplate(&tpl)

	if err := act.Save(); err != nil {
		fmt.Printf("failed to create action: %v\n", err)
		return
	}

	fmt.Println("succeed creating action:")
	o, _ := act.JsonPretty()
	fmt.Println(string(o))
}

func createCase(defJsonFile string) {
	defjson, err := ioutil.ReadFile(defJsonFile)
	if err != nil {
		fmt.Printf("failed to read definition file of case: %v\n", err)
		return
	}

	ca, err := TC.FromJson(defjson)
	if err != nil {
		fmt.Printf("failed to parse definition of case: %v\n", err)
		return
	}

	if err := ca.Save(); err != nil {
		fmt.Printf("failed to create case: %v\n", err)
		return
	}
	fmt.Println("succeed creating case:")
	o, _ := ca.JsonPretty()
	fmt.Println(string(o))
}
