package main

import (
	"fmt"
	"github.com/fzzy/radix/extra/pool"
	"os"
	"strings"

	C "legolas/common/config"
	A "legolas/common/models/action"
	J "legolas/common/models/job"
	R "legolas/common/models/run"
	T "legolas/common/models/template"
	TC "legolas/common/models/testcase"
	S "legolas/common/storage"
)

func main() {
	if len(os.Args) < 2 {
		showHelp()
		os.Exit(0)
	}

	redisPool, err := pool.NewPool("tcp", C.RedisHost, C.RedisPoolSize)
	if err != nil {
		fmt.Printf("Cannot create Redis pool at: %s\n", C.RedisHost)
	}
	J.SetRedisPool(redisPool)

	mongo := S.AskForMongo()
	defer mongo.Close()

	T.SetCol(mongo)
	TC.SetCol(mongo)
	A.SetCol(mongo)
	R.SetCol(mongo)

	subCommand := strings.ToLower(os.Args[1])
	switch subCommand {
	case "help":
		if len(os.Args) > 2 {
			helpCommand(os.Args[2])
		} else {
			showHelp()
		}
	case "create", "new":
		createCommand(os.Args)
	case "show":
		showCommand(os.Args)
	case "list":
		listCommand(os.Args)
	case "run":
		runCommand(os.Args)
	default:
		fmt.Println("unsupported sub commands " + subCommand)
	}
}

func showHelp() {
	text := `
	Legolas Cli tool v0.1

	Sub commands: create, list, show, run
	Use 'legolas-cli help <subcommand>' for details
	`
	fmt.Println(text)
}

func helpCommand(subCommand string) {
	var text string

	switch subCommand {
	case "show":
		text = `show details of data entries`
	default:
		text = "unknown subcommand: " + subCommand
	}

	fmt.Println(text)
}
