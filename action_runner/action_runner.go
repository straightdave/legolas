package main

import (
	"flag"
	"fmt"
	"github.com/fzzy/radix/redis"
	"log"
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

func main() {
	flag.Parse()

	// set log to local file
	logf, err := os.OpenFile(*logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Cannot access to log file: %s\n", *logFile)
	}
	defer logf.Close()
	log.SetOutput(logf)

	client, err := redis.Dial("tcp", *queueAddress)
	if err != nil {
		log.Fatalf("Cannot connect to Redis instance: %s\n", *queueAddress)
	}
	defer client.Close()

	// handling ctrl-c/TERM, to exit gracefully
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		logf.Close()
		client.Close()
		fmt.Println("legolas action runner is stoppeed.")
	}()

	fmt.Println("action runner is listening...")
	log.Println("action runner stared at: ", time.Now().String())
	for {
		data, err := client.Cmd("BLPOP", *queueName, 0).List()
		if err != nil {
			if _, ok := err.(*redis.CmdError); ok {
				continue
			} else {
				log.Fatalf("BLPOP met fatal error: %v\n", err)
			}
		}

		if len(data) > 1 {
			go jobHandler(data[1])
		}
	}
}
