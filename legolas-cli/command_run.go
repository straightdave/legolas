package main

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"time"

	A "legolas/common/models/action"
	J "legolas/common/models/job"
	R "legolas/common/models/run"
	TC "legolas/common/models/testcase"
)

func runCommand(args []string) {
	if len(args) < 3 {
		fmt.Println("...which case? (case path + name)")
		return
	}

	if len(args) == 3 {
		fmt.Println("currently not support run all cases under path")
		return
	}

	cpath := args[2]
	cname := args[3]

	tc, err := TC.GetOne(cpath, cname)
	if err != nil {
		fmt.Printf("failed to get case by %s/%s: %v\n", cpath, cname, err)
		return
	}

	run := R.New(tc.Id)
	run.StartedAt = time.Now()
	if err := run.Save(); err != nil {
		fmt.Printf("failed to save run: %v\n", err)
		return
	}

	actions, err := A.GetAllByCaseId(tc.Id)
	if err != nil {
		fmt.Printf("failed to get actions: %v\n", err)
		return
	}

	var prev bson.ObjectId
	for _, act := range actions {
		job := &J.Job{
			RunId:        run.Id,
			ActionId:     act.Id,
			PrevActionId: prev,
		}

		if err := J.Append(job); err != nil {
			fmt.Printf("failed to push job into queue: %v\n", err)
			break
		}

		prev = act.Id
	}
	fmt.Println("successfully created the test run: " + run.Id)
}
