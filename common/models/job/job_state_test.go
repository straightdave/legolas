package job

import (
	"fmt"
	"testing"
	"time"
)

func TestSaveJobState(t *testing.T) {
	js := NewJobState("case_run_001", "action-1", "$case/path", "case-1")
	js.AddResult("result-1", 123)
	js.AddResult("result-2", "string result")
	js.StartedAt = time.Now()
	if err := js.Save(); err != nil {
		t.Fatalf("cannot save job state: %v\n", err)
	}

	js2, err := GetJobState("case_run_001", "action-1")
	if err != nil {
		t.Fatalf("cannot get job state: %v\n", err)
	}

	c, _ := js2.JsonPretty()
	fmt.Println(string(c))
}
