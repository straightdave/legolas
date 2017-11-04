/*
 */

package main

import (
	"legolas/server"
	"log"
)

func main() {
	var _s server.Server
	if err := _s.Start(); err != nil {
		log.Fatal("Err: %v", err)
	}
}
