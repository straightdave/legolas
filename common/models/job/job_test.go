package job

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestJobInJson(t *testing.T) {
	data1 := `{"case_run_id":"case_run_1","name":"action1","snippet":"print(\"hello python\")","pre_action":"case_0"}`
	data2 := Job{
		CaseRunID: "case_run_1",
		Name:      "action1",
		Snippet:   `print("hello python")`,
		PreAction: "action0",
	}

	c_job, err := JobFromJson([]byte(data1))
	if err != nil {
		t.Errorf("failed json unmarshalling\n")
	}

	fmt.Println(c_job.CaseRunID)
	fmt.Println(c_job.Name)
	fmt.Println(c_job.Snippet)
	fmt.Println(c_job.PreAction)

	data2_j, err := data2.Json()
	if err != nil {
		t.Errorf("failed json marshalling\n")
	}
	fmt.Println(string(data2_j))
	fmt.Println("-----------------")
}

func TestJobInBson(t *testing.T) {
	var jobs = []Job{
		{
			CaseRunID: "case_run_1",
			Name:      "action1",
			Snippet:   `print("hello python!")`,
			PreAction: "case_run_0",
		},
	}

	for _, job := range jobs {
		data, err := job.Json()
		if err != nil {
			t.Errorf("marshal failed: %v\n", err)
		}
		fmt.Println(string(data))

		// test the snippet could be run after unmarshalling
		c_job, err := JobFromJson(data)
		if err != nil {
			t.Errorf("content cannot unmarshal\n")
		}

		tempFileName := fmt.Sprintf("%s_snippet.py", c_job.Name)
		err = ioutil.WriteFile(tempFileName, []byte(strings.Trim(string(c_job.Snippet), "\n ")), 0755)
		if err != nil {
			t.Errorf("cannopt write snippet to file: %v\n", err)
		}
		defer func() {
			fmt.Println("remove tempfile:", tempFileName)
			err := os.Remove(tempFileName)
			if err != nil {
				t.Fatalf("cannot remove temp file: %s: %v\n", tempFileName, err)
			}
		}()

		cmd := exec.Command("python", tempFileName)
		cmdOut, err := cmd.Output()
		if err != nil {
			t.Errorf("cannot get python process output: %v\n", err)
		} else {
			fmt.Printf("%s : %s\n", c_job.Name, string(cmdOut))
		}

		fmt.Println("----------------")
	}
}
