package main

import (
	"flag"
	"fmt"
	"github.com/fzzy/radix/extra/pool"
	L "log"
	"os"
	"os/signal"
	"syscall"

	"legolas/common/config"
	"legolas/common/models/job"
)

var (
	queueAddress = flag.String("q", config.RedisHost, "Redis instance address")
	queueName    = flag.String("n", config.Queue, "Redis queue name")
	logFile      = flag.String("l", config.RunnerLogFile, "Log file name")
)

func main() {
	flag.Parse()

	// set logger
	logf, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Cannot access to log file: %s\n", *logFile)
		os.Exit(1)
	}
	L.SetOutput(logf)
	L.Println()

	// set redis pool
	redisPool, err := pool.NewPool("tcp", *queueAddress, config.RedisPoolSize)
	if err != nil {
		L.Fatalf("Cannot create Redis pool at: %s\n", *queueAddress)
	}
	job.SetRedisPool(redisPool)

	// set deferred cleanup
	defer func() {
		logf.Close()
		redisPool.Empty()
	}()

	// handling ctrl-c and TERM signals gracefully
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		L.Fatalln("legolas runner is interrupted.")
	}()

	L.Println("legolas runner started")
	for {
		j, err := job.Pop()
		if err != nil {
			L.Printf("failed to get job from the queue: %v\n", err)
			continue
		}
		go handle(&j)
	}
}
