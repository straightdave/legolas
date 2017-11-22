package helpers

import (
	"fmt"
	"io/ioutil"
	"strings"
)

var harness_before = `
import sys
import time
from legolas import Legolas
`

var harness_after = `
if __name__ == "__main__":
    ctx = Legolas(sys.argv[1])
    ctx.save_result('started_at', time.ctime())
    action_main(ctx)
    ctx.save_result('ended_at', time.ctime())
`

func GenScript(fileName, snippet string) error {
	text := fmt.Sprintf("%s\n%s\n%s", harness_before, snippet, harness_after)
	return ioutil.WriteFile(fileName, []byte(strings.Trim(text, "\n ")), 0644)
}
