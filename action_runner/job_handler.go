package main

import (
	"fmt"
	"io/ioutil"
	M "legolas/common/models"
	L "log"
	"os/exec"
	"strings"
)

func handle(job *M.Job) {
	L.Printf("Get job: [\n%v\n]\n", *job)

	r, err := redisPool.Get()
	if err != nil {
		L.Printf("cannot get a redis connection from pool: %v\n", err)
		return
	}
	defer redisPool.Put(r)

	js := &M.JobStateProxy{Id: job.JobID(), Rc: r}

	// if any panic happens in this func, mark it failed.
	defer func() {
		if p := recover(); p != nil {
			js.Fail(p.(error).Error())
		}
	}()

	// check previous action state
	if job.PreAction != "" {
		jsPrev := &M.JobStateProxy{Id: job.PrevJobID(), Rc: r}
		switch jsPrev.State() {
		case "done": // go on
			break
		case "failed": // just discard and mark as failed
			js.Fail(fmt.Sprintf("Aborted due to previous Job [%s] was failed.", jsPrev.Id))
			return
		default: // postpone
			L.Printf("previous job [%s] is not done yet. append to queue.\n", jsPrev.Id)
			jp := &M.JobProxy{Queue: *queueName, Rc: r}
			jp.Append(job)
			return
		}
	}

	// mark as in-progress. discard job if this fails.
	if err := js.InProgress(); err != nil {
		L.Printf("failed to set job as in-progress: %v\n", err)
		return
	}

	// save snippet to local temp file
	fn := fmt.Sprintf("%s__%s.py", job.JobID(), RandSuffix4())
	err = ioutil.WriteFile(fn, []byte(strings.Trim(job.Snippet, "\n ")), 0755)
	if err != nil {
		L.Printf("cannot create snippet file: %v\n", err)
		js.Fail(err.Error())
		return
	}
	// defer os.Remove(fn)

	// execute snippet, collecting outputs
	cmd := exec.Command("python", fn)
	cmdOut, err := cmd.CombinedOutput()
	if err != nil {
		L.Printf("cannot run and get output of process: %v\n", err)
		js.Fail(err.Error())
		return
	}

	js.Done()
	js.AppendLog(string(cmdOut))
}
