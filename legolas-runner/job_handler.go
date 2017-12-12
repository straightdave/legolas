package main

import (
	"fmt"
	L "log"
	"os/exec"
	"time"

	C "legolas/common/config"
	H "legolas/common/helpers"
	A "legolas/common/models/action"
	J "legolas/common/models/job"
	E "legolas/common/models/jobstate"
	R "legolas/common/models/run"
	S "legolas/common/storage"
)

func handle(job *J.Job) {
	ch := make(chan bool, 1)
	go func() {
		handling(job)
		ch <- true
	}()

	select {
	case <-ch:
		L.Printf("Job handling completed in time\n")

	case <-timeout(C.JobTimeout):
		L.Printf("Job processing timeout (%d seconds max)\n", C.JobTimeout)
		setJobAsTimeout(job)

	default:
	}
}

func handling(job *J.Job) {
	jid := job.Id()
	L.Printf("[%s] Handler: got job\n", jid)

	mongo := S.AskForMongo()
	defer mongo.Close()

	A.SetCol(mongo)
	E.SetCol(mongo)
	R.SetCol(mongo)

	// get its run
	r, err := R.GetOne(job.RunId)
	if err != nil {
		L.Printf("[%s] failed to get run data: %v\n", jid, err)
		return
	}

	// get its action
	act, err := A.GetOneById(job.ActionId)
	if err != nil {
		L.Printf("[%s] failed to get action data: %v\n", jid, err)
		return
	}

	// create a job state for it
	js := E.New(job.RunId, job.ActionId)
	js.State = C.Running
	if err := js.Save(); err != nil {
		L.Printf("[%s] failed to create job state: %v\n", jid, err)
		return
	}

	// modify its run's output to inprogress
	r.Output = js.State
	if err := r.Save(); err != nil {
		L.Printf("[%s] failed to save run modification: %v\n", jid, err)
		return
	}

	defer func() {
		if p := recover(); p != nil {
			js.State = C.Failed
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
		case C.Done: // continue
			break
		case C.Failed, C.Aborted, C.Timeout: // just discard it and mark as aborted
			js.State = C.Aborted
			js.Error = fmt.Sprintf("previous job [action: %s] is [%s]", job.PrevActionId, prev.State)
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

	if act.IsMocking {
		// if the action is a mockingbird,
		// directly save all mock data into results in jobstate

		L.Printf("[%s] it's mockingbird. Saving mock data as results.", jid)

		js, err := E.GetOne(job.RunId, job.ActionId)
		if err != nil {
			L.Printf("[%s] failed to get job state: %v\n", jid, err)
			return
		}

		js.Results = act.MockData
		js.State = C.Done
		js.EndedAt = time.Now().Format(time.ANSIC)
		if err := js.Save(); err != nil {
			L.Printf("[%s] failed to save job state: %v\n", jid, err)
			return
		}

		// get its run
		run, err := R.GetOne(job.RunId)
		if err != nil {
			L.Printf("[%s] failed to get run data: %v\n", jid, err)
			return
		}

		// modify run's output and end time
		run.Output = js.State
		run.EndedAt = time.Now()
		if err := run.Save(); err != nil {
			L.Printf("[%s] failed to save run modification: %v\n", jid, err)
			return
		}
	}

	// processing the action
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

	// create context data for script execution
	ctx, err := job.Json()
	if err != nil {
		L.Printf("[%s] failed to serialize job as script context (json): %v\n", jid, err)
		return
	}
	ctxStr := string(ctx)
	cmd := exec.Command("python", fn, ctxStr)
	cmdOut, _ := cmd.CombinedOutput()
	L.Printf("[%s] >>>\n%s\n", jid, cmdOut)

	// finishing the job run, update some fields
	// re-fetch the job state for updating the latest
	// TODO: provide 'refresh' methods for those models
	js2, err := E.GetOne(job.RunId, job.ActionId)
	if err != nil {
		L.Printf("[%s] cannot re-fetch job state: %v\n", err)
		return
	}

	js2.Output = string(cmdOut)
	if js2.State != C.Failed && js2.State != C.Aborted {
		js2.State = C.Done
	}

	if err := js2.Save(); err != nil {
		L.Printf("[%s] failed to set job state to done: %v\n", jid, err)
		return
	}

	// get its run again
	r2, err := R.GetOne(job.RunId)
	if err != nil {
		L.Printf("[%s] failed to get run data: %v\n", jid, err)
		return
	}

	// modify its run's output and end time
	r2.Output = js2.State
	r2.EndedAt = time.Now()
	if err := r2.Save(); err != nil {
		L.Printf("[%s] failed to save run modification: %v\n", jid, err)
	}
}

func setJobAsTimeout(job *J.Job) {
	mongo := S.AskForMongo()
	defer mongo.Close()
	E.SetCol(mongo)

	if js, err := E.GetOne(job.RunId, job.ActionId); err == nil {
		js.State = C.Timeout
		if err := js.Save(); err != nil {
			L.Printf("Failed to save job state: %v\n", err)
		}
	}
}

func timeout(seconds int) <-chan time.Time {
	return time.After(time.Second * time.Duration(seconds))
}
