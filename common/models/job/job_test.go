package job

import (
	"fmt"
	"github.com/fzzy/radix/extra/pool"
	"testing"

	"legolas/common/config"
)

func TestCreateJob(t *testing.T) {
	p, err := pool.NewPool("tcp", config.RedisHost, 3)
	if err != nil {
		t.Fatalf("cannot init redis pool: %v\n", err)
	}
	SetRedisPool(p)

	job := &Job{
		CaseRunID:  "test-run-id",
		CasePath:   "$case/path",
		CaseName:   "case-1",
		ActionName: "action-1",
	}

	c, err := job.JsonPretty()
	if err != nil {
		t.Fatalf("cannot parse popped job: %v\n", err)
	}
	fmt.Println(string(c))

	if err := Append(job); err != nil {
		t.Fatalf("cannot push new job: %v\n", err)
	}

	job2, err := Pop()
	if err != nil {
		t.Fatalf("cannot pop job: %v\n", err)
	}

	c, err = job2.JsonPretty()
	if err != nil {
		t.Fatalf("cannot parse popped job: %v\n", err)
	}
	fmt.Println(string(c))
}
