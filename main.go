/*
 */

package main

import (
	"fmt"
	"legolas/server"
	"os"
)

func main() {
	var _server server.Server
	if err := _server.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Err: %v", err)
	}
}
