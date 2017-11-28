package main

import (
	"fmt"
	L "log"
	"os/exec"

	C "legolas/common/config"
	H "legolas/common/helpers"
	A "legolas/common/models/action"
	J "legolas/common/models/job"
	E "legolas/common/models/jobstate"
	S "legolas/common/storage"
)

func handle(job *J.Job) {
	// for the sake of not breaking the runner,
	// eat all panics here
	defer func() {
		if p := recover(); p != nil {
			L.Printf("error occured in one job handling process: %v\n", p.(error))
		}
	}()

	jid := job.Id()
	L.Printf("[%s] Handler: got job\n", jid)

	mongo := S.AskForMongo()
	defer mongo.Close()

	A.SetCol(mongo)
	E.SetCol(mongo)

	act, err := A.GetOneById(job.ActionId)
	if err != nil {
		L.Printf("[%s] failed to get action data: %v\n", jid, err)
		return
	}

	js := E.New(job.RunId, job.ActionId)
	js.State = "running"
	if err := js.Save(); err != nil {
		L.Printf("[%s] failed to create job state: %v\n", jid, err)
		return
	}

	defer func() {
		if p := recover(); p != nil {
			js.State = "failed"
			js.Error = p.(error).Error()
			if err := js.Save(); err != nil {
				L.Printf("[%s] panic in panic: failed to set job state as failed: %v\n", jid, err)
			}
		}
	}()

	// check previous job state
	if job.PrevActionId != "" {
		prev, err := E.GetOne(job.RunId, job.PrevActionId)
		if err != nil {
			L.Printf("[%s] failed to get previous job state: %v\n", jid, err)
			return
		}

		switch prev.State {
		case "done": // continue
			break
		case "failed", "aborted": // just discard it and mark as aborted
			js.State = "aborted"
			js.Error = fmt.Sprintf("previous job [action: %s] is failed/aborted", job.PrevActionId)
			if err := js.Save(); err != nil {
				L.Printf("[%s] failed to set job state: %v\n", jid, err)
			}
			return
		default: // postpone it
			L.Printf("[%s] previous job [action: %s] is not done yet. append to queue.\n", jid, job.PrevActionId)
			if err := J.Append(job); err != nil {
				L.Printf("[%s] failed to postpone job which will be discarded: %v\n", jid, err)
			}
			return
		}
	}

	// generate script file
	fn := fmt.Sprintf("%s/%s__%s__%s.py", C.ScriptHive, job.RunId.Hex(), job.ActionId.Hex(), H.RandSuffix4())
	snippet, err := act.Snippet()
	if err != nil {
		L.Printf("[%s] failed to get snippet: %v\n", err)
		return
	}
	if err := H.GenScript(fn, snippet); err != nil {
		L.Printf("[%s] cannot create script file [%s]: %v\n", jid, fn, err)
		return
	}

	// make context data for script execution
	ctx, err := job.Json()
	if err != nil {
		L.Printf("[%s] failed to serialize job as script context (json): %v\n", jid, err)
		return
	}
	ctxStr := string(ctx)
	cmd := exec.Command("python", fn, ctxStr)
	cmdOut, _ := cmd.CombinedOutput()
	L.Printf("[%s] >>>\n%s\n", jid, cmdOut)

	// complete job run
	// after python process done, re-fetch the job state for updating the latest
	js2, err := E.GetOne(job.RunId, job.ActionId)
	if err != nil {
		L.Printf("[%s] cannot re-fetch job state: %v\n", err)
		return
	}
	js2.Output = string(cmdOut)
	if js2.State != "failed" && js2.State != "aborted" {
		js2.State = "done"
	}
	if err := js2.Save(); err != nil {
		L.Printf("[%s] failed to set job state to done: %v\n", jid, err)
	}
}
