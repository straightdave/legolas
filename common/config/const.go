/*
   shared constants
*/
package config

const (
	MongoHost     string = "localhost"
	RedisHost     string = "localhost:6379"
	RedisPoolSize int    = 10
	Queue         string = "actionjobs"
	RunIdLength   int    = 8
	RunnerLogFile string = "legolas_runner.log"
	ScriptHive    string = "/Users/wei.wu/myrepos/src/legolas/legolas-runner/hive"
)

const (
	NotStarted string = "notstarted"
	Running    string = "running"
	Failed     string = "failed"
	Aborted    string = "aborted"
	Done       string = "done"
)
