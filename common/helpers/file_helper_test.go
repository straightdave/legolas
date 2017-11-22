package helpers

import (
	"testing"
)

func TestGenScript(t *testing.T) {
	snippet := `print("hello world")`

	if err := GenScript("temp.py", snippet); err != nil {
		t.Fatalf("cannot gen script: %v\n", err)
	}
}
