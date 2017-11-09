package main

import (
	"fmt"
	"io/ioutil"
	"legolas/common/models"
	"log"
	"os"
	"os/exec"
	"strings"
)

func jobHandler(data string) {
	log.Printf("Get => [ %s ]", data)

	job, err := models.JobFromJson([]byte(data))
	if err != nil {
		log.Printf("Unmarshalling failed: %v\n", err)
	}

	log.Printf("snippet => {\n%s\n}\n", job.Snippet)

	tempFileName := fmt.Sprintf("%s#%s_snippet.py", job.CaseRunID, job.Name)
	err = ioutil.WriteFile(tempFileName, []byte(strings.Trim(job.Snippet, "\n ")), 0755)
	if err != nil {
		log.Printf("cannopt write snippet to file: %v\n", err)
	}
	defer func() {
		log.Println("remove tempfile:", tempFileName)
		err := os.Remove(tempFileName)
		if err != nil {
			log.Printf("cannot remove temp file: %s: %v\n", tempFileName, err)
		}
	}()

	cmd := exec.Command("python", tempFileName)
	cmdOut, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("cannot get output of process: %v\n", err)
	}

	log.Printf("output => {\n%s\n}\n", strings.Trim(string(cmdOut), "\n "))
}
