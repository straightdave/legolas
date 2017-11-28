package main

import (
	"fmt"

	A "legolas/common/models/action"
	R "legolas/common/models/run"
	T "legolas/common/models/template"
	TC "legolas/common/models/testcase"
)

func listCommand(args []string) {
	if len(args) < 3 {
		fmt.Println("...what? (cases, actions, templates, runs)")
		return
	}

	switch args[2] {
	case "runs":
		cpath := args[3]
		cname := args[4]

		ca, err := TC.GetOne(cpath, cname)
		if err != nil {
			fmt.Printf("failed to get case of %s/%s: %v\n", cpath, cname, err)
		}

		runs, err := R.GetAll(&ca)
		if err != nil {
			fmt.Printf("failed to list runs of cpath:[%s], cname:[%s]: %v\n", cpath, cname, err)
			break
		}

		if len(runs) < 1 {
			fmt.Printf("no runs found for case -> cpath:[%s], cname:[%s]\n", cpath, cname)
			break
		}

		for _, run := range runs {
			o, _ := run.JsonPretty()
			fmt.Println(string(o))
		}

	case "cases":
		path := args[3]

		cases, err := TC.GetAll(path)
		if err != nil {
			fmt.Printf("failed to list cases of path:[%s]: %v\n", path, err)
			break
		}

		if len(cases) < 1 {
			fmt.Printf("no cases found under path:[%s]\n", path)
			break
		}

		for _, ca := range cases {
			o, _ := ca.JsonPretty()
			fmt.Println(string(o))
		}

	case "actions":
		cpath := args[3]
		cname := args[4]

		acts, err := A.GetAll(cpath, cname)
		if err != nil {
			fmt.Printf("failed to list actions of cpath:[%s], cname:[%s]: %v\n", cpath, cname, err)
			break
		}

		if len(acts) < 1 {
			fmt.Printf("no cases found under cpath:[%s], cname:[%s]\n", cpath, cname)
			break
		}

		for _, act := range acts {
			o, _ := act.JsonPretty()
			fmt.Println(string(o))
		}

	case "templates":
		path := args[3]

		tpls, err := T.GetAll(path)
		if err != nil {
			fmt.Printf("failed to list templates of path:[%s]: %v\n", path, err)
			break
		}

		if len(tpls) < 1 {
			fmt.Printf("no template found under path:[%s]\n", path)
			break
		}

		for _, tpl := range tpls {
			o, _ := tpl.JsonPretty()
			fmt.Println(string(o))
		}

	default:
		fmt.Println("unknown object type " + args[2])
	}
}
