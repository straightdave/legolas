package main

import (
	"fmt"
	"github.com/fzzy/radix/extra/pool"
	L "log"
	"os"
	"os/signal"
	"syscall"

	C "legolas/common/config"
	J "legolas/common/models/job"
)

func main() {
	f, err := os.OpenFile(C.RunnerLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Cannot access to log file: %s\n", C.RunnerLogFile)
		os.Exit(1)
	}
	L.SetOutput(f)

	redisPool, err := pool.NewPool("tcp", C.RedisHost, C.RedisPoolSize)
	if err != nil {
		L.Printf("Cannot initialize redis pool: %v\n", err)
		f.Close()
		os.Exit(1)
	}
	J.SetRedisPool(redisPool)

	cleanUp := func() {
		f.Close()
		redisPool.Empty()
	}
	defer cleanUp()

	// handling ctrl-c and TERM signals gracefully
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		L.Println("legolas runner is interrupted.")
		cleanUp()
		os.Exit(1)
	}()

	L.Println("legolas runner started")
	for {
		job, err := J.Pop()
		if err != nil {
			L.Printf("failed to get job from the queue: %v\n", err)
			continue
		}
		go handle(&job)
	}
}
