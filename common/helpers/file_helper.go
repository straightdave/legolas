package helpers

import (
	"fmt"
	"io/ioutil"
	"strings"
)

var harness_before = `
import sys
from lib import legolas
`

var harness_after = `
if __name__ == "__main__":
    ctx = legolas.Legolas(sys.argv[1])
    ctx._set_start_time()

    try:
        action_main(ctx)
        ctx._upload_results()
    except:
        print('Err: ' + str(sys.exc_info()))
    else:
        ctx._set_end_time()
`

func GenScript(fileName, snippet string) error {
	text := fmt.Sprintf("%s\n%s\n%s", harness_before, snippet, harness_after)
	return ioutil.WriteFile(fileName, []byte(strings.Trim(text, "\n ")), 0644)
}
