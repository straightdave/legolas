package main

import (
	"flag"
	"fmt"
	"github.com/fzzy/radix/extra/pool"
	"legolas/common/models"
	L "log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	queueAddress = flag.String("q", "localhost:6379", "Redis instance address")
	queueName    = flag.String("n", "actionjobs", "Redis queue name")
	logFile      = flag.String("l", "action_runner.log", "Log file name")
)

var (
	poolSize  = 10
	redisPool *pool.Pool
)

func main() {
	flag.Parse()

	logf, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Cannot access to log file: %s\n", *logFile)
		os.Exit(1)
	}
	L.SetOutput(logf)

	redisPool, err = pool.NewPool("tcp", *queueAddress, poolSize)
	if err != nil {
		L.Fatalf("Cannot create Redis pool at: %s\n", *queueAddress)
	}

	rc, err := redisPool.Get()
	if err != nil {
		L.Fatalf("Cannot get a connection from Redis pool for the main goroutine: %v\n", err)
	}

	jobProxy := &models.JobProxy{Queue: *queueName, Rc: rc}

	cleanup := func() {
		logf.Close()
		redisPool.Put(rc)
		redisPool.Empty()
	}
	defer cleanup()

	// handling ctrl-c and TERM signals
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		fmt.Println("legolas action runner is interrupted.")
		os.Exit(1)
	}()

	fmt.Println("action runner is listening...")
	L.Println("action runner stared at: ", time.Now().String())
	for {
		job, err := jobProxy.Pop()
		if err != nil {
			L.Printf("failed to get job from the queue: %v\n", err)
			continue
		}
		go handle(job)
	}
}
