package main

import (
	"fmt"
	L "log"
	"os/exec"

	"legolas/common/config"
	"legolas/common/helpers"
	A "legolas/common/models/action"
	J "legolas/common/models/job"
)

func handle(job *J.Job) {
	jid := job.Id()
	L.Printf("Handler: get job id:[%s]\n", jid)

	// get action data by job info
	act, err := A.GetAction(job.CasePath, job.CaseName, job.ActionName)
	if err != nil {
		L.Printf("[%s] failed to get action data: %v\n", jid, err)
		return
	}

	// create job state
	js := J.NewJobState(job.CaseRunID, job.ActionName, job.CasePath, job.CaseName)
	if err := js.Save(); err != nil {
		L.Printf("[%s] failed to create job state: %v\n", jid, err)
		return
	}

	// if any panic happens in this func, mark it failed.
	defer func() {
		if p := recover(); p != nil {
			js.State = "failed"
			js.Error = p.(error).Error()
			if err := js.Save(); err != nil {
				L.Printf("[%s] panic in panic: failed to set job state as failed: %v\n", jid, err)
			}
		}
	}()

	// check previous action
	if job.PrevAction != "" {
		prevActionState, err := J.GetJobState(job.CaseRunID, job.PrevAction)
		if err != nil {
			L.Printf("[%s] failed to get previous job state: %v\n", jid, err)
			return
		}

		switch prevActionState.State {
		case "done": // continue
			break
		case "failed", "aborted": // just discard and mark as aborted
			js.State = "aborted"
			js.Error = fmt.Sprintf("previous job [%s] is in the state failed/aborted", job.PrevAction)
			if err := js.Save(); err != nil {
				L.Printf("[%s] failed to set job state: %v\n", jid, err)
			}
			return
		default: // postpone
			L.Printf("[%s] previous job [%s] is not done yet. append to queue.\n", jid, job.PrevAction)
			if err := J.Append(job); err != nil {
				L.Printf("[%s] failed to postpone job: %v\n", jid, err)
			}
			return
		}
	}

	// save snippet to local file
	// the file name would be <run_id>__<action name>__<random suffix>.py
	fn := fmt.Sprintf("%s/%s__%s__%s.py", config.ScriptHive, job.CaseRunID, job.ActionName, helpers.RandSuffix4())
	if err := helpers.GenScript(fn, act.Snippet); err != nil {
		L.Printf("[%s] cannot create snippet file [%s]: %v\n", jid, fn, err)
		return
	}

	// execute script, collecting outputs
	// context is job (in json)
	ctx, err := job.Json()
	if err != nil {
		L.Printf("[%s] failed to serialize job as script context (json): %v\n", jid, err)
		return
	}

	ctxStr := string(ctx)
	L.Printf("[%s] run script:[%s] with ctx:[%s]\n", jid, fn, ctxStr)

	cmd := exec.Command("python", fn, ctxStr)
	cmdOut, err := cmd.CombinedOutput()
	if err != nil {
		L.Printf("[%s] cannot run and get output of process: %v\n", jid, err)
		return
	}
	L.Println(">>> " + string(cmdOut))

	// complete job run
	js.Output = string(cmdOut)
	js.State = "done"
	if err := js.Save(); err != nil {
		L.Printf("[%s] failed to set job state to done: %v\n", jid, err)
	}
}
