package main

import (
	"fmt"
	"github.com/fzzy/radix/extra/pool"
	"io/ioutil"
	"os"
	"strings"

	"legolas/common/config"
	"legolas/common/models/action"
	"legolas/common/models/job"
	"legolas/common/models/run"
	"legolas/common/models/template"
	"legolas/common/models/testcase"
)

func main() {
	if len(os.Args) < 2 {
		showHelp()
		os.Exit(0)
	}

	// set redis pool
	redisPool, err := pool.NewPool("tcp", config.RedisHost, config.RedisPoolSize)
	if err != nil {
		fmt.Printf("Cannot create Redis pool at: %s\n", config.RedisHost)
	}
	job.SetRedisPool(redisPool)

	subCommand := strings.ToLower(os.Args[1])
	switch subCommand {
	case "help":
		showHelp()
	case "create", "new":
		createCommand(os.Args)
	case "list", "show":
		listCommand(os.Args)
	case "run":
		runCommand(os.Args)
	default:
		fmt.Println("unsupported sub commands " + subCommand)
	}
}

func showHelp() {
	fmt.Println("Legolas Cli v0.1 - dave")
	fmt.Println("subcommands: help, create/new, list/show, run")
}

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
		os.Exit(1)
	}

	tpl, err := template.FromJson(defjson)
	if err != nil {
		fmt.Printf("failed to parse definition of template: %v\n", err)
		os.Exit(1)
	}

	snippet, err := ioutil.ReadFile(snippetFile)
	if err != nil {
		fmt.Printf("failed to read snippet file of template: %v\n", err)
		os.Exit(1)
	}

	tpl.Snippet = string(snippet)
	if err := tpl.Save(); err != nil {
		fmt.Printf("failed to create new template: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("succeed creating template:")
	o, _ := tpl.JsonPretty()
	fmt.Println(string(o))
}

func createAction(defJsonFile string) {
	defjson, err := ioutil.ReadFile(defJsonFile)
	if err != nil {
		fmt.Printf("failed to read definition file of action: %v\n", err)
		os.Exit(1)
	}

	act, err := action.FromJson(defjson)
	if err != nil {
		fmt.Printf("failed to parse definition of action: %v\n", err)
		os.Exit(1)
	}

	err = act.ApplyTemplate(act.TemplatePath, act.TemplateName)
	if err != nil {
		fmt.Printf("failed to apply template [%s/%s]: %v\n", act.TemplatePath, act.TemplateName, err)
		os.Exit(1)
	}

	if err := act.Save(); err != nil {
		fmt.Printf("failed to create action: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("succeed creating action:")
	o, _ := act.JsonPretty()
	fmt.Println(string(o))
}

func createCase(defJsonFile string) {
	defjson, err := ioutil.ReadFile(defJsonFile)
	if err != nil {
		fmt.Printf("failed to read definition file of case: %v\n", err)
		os.Exit(1)
	}

	ca, err := testcase.FromJson(defjson)
	if err != nil {
		fmt.Printf("failed to parse definition of case: %v\n", err)
		os.Exit(1)
	}

	if err := ca.Save(); err != nil {
		fmt.Printf("failed to create case: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("succeed creating case:")
	o, _ := ca.JsonPretty()
	fmt.Println(string(o))
}

func runCommand(args []string) {
	if len(args) < 3 {
		fmt.Println("...which case? (case path + name)")
		return
	}

	var cpath, cname string
	if len(args) == 3 {
		fmt.Println("currently not support run all cases under path")
		return
	} else if len(args) > 3 {
		cpath = args[2]
		cname = args[3]
	}

	r, err := run.NewRun(cpath, cname)
	if err != nil {
		fmt.Printf("failed to run: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("succeed creating run:")
	o, _ := r.JsonPretty()
	fmt.Println(string(o))
}

func listCommand(args []string) {
	if len(args) < 3 {
		fmt.Println("...what? (cases, case, actions, action, templates, template, runs, run, jobs, job)")
		return
	}

	switch args[2] {
	case "jobs":
		jobs, err := job.View()
		if err != nil {
			fmt.Printf("failed to list jobs: %v\n", err)
			break
		}

		if len(jobs) < 1 {
			fmt.Printf("no jobs found\n")
			break
		}

		for _, j := range jobs {
			o, _ := j.JsonPretty()
			fmt.Println(string(o))
		}

	case "runs":
		cpath := args[3]
		cname := args[4]
		runs, err := run.GetRuns(cpath, cname)
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
		cases, err := testcase.GetTestCases(path)
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

	case "case":
		path := args[3]
		name := args[4]
		c, err := testcase.GetTestCase(path, name)
		if err != nil {
			fmt.Printf("failed to list cases of path:[%s], name:[%s]: %v\n", path, name, err)
			break
		}
		o, _ := c.JsonPretty()
		fmt.Println(string(o))

	case "actions":
		cpath := args[3]
		cname := args[4]
		acts, err := action.GetActions(cpath, cname)
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

	case "action":
		cpath := args[3]
		cname := args[4]
		name := args[5]
		act, err := action.GetAction(cpath, cname, name)
		if err != nil {
			fmt.Printf("failed to get action of cpath:[%s], cname:[%s], name:[%s]: %v\n", cpath, cname, name, err)
			break
		}
		o, _ := act.JsonPretty()
		fmt.Println(string(o))

	case "templates":
		path := args[3]
		tpls, err := template.GetTemplates(path)
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

	case "template":
		path := args[3]
		name := args[4]
		tpl, err := template.GetTemplate(path, name)
		if err != nil {
			fmt.Printf("failed to list templates of path:[%s], name:[%s]: %v\n", path, name, err)
			break
		}
		o, _ := tpl.JsonPretty()
		fmt.Println(string(o))

	default:
		fmt.Println("unknown object type " + args[2])
	}
}
