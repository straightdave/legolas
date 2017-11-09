package models

import (
	"fmt"
	"testing"
)

func TestCaseRunState(t *testing.T) {
	var states = []CaseRunState{
		{
			CaseRunId: "case_run_1",
			Context: map[string]interface{}{
				"cxt1": "a string",
				"cxt2": map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
				"cxt3": []byte(`wahahah wahahaha hahah`),
			},
		},
	}

	for _, state := range states {
		// test marshalling
		data, err := state.ToJson()
		if err != nil {
			t.Errorf("failed to marshal: %v\n", err)
		}
		fmt.Println(string(data))

		// test unmarshalling
		c_state, err := CaseRunStateFromJson(data)
		if err != nil {
			t.Errorf("cannot unmarshal: %v\n", err)
		}

		fmt.Println("cxt1: ", c_state.Context["cxt1"])
		fmt.Println("cxt2: ", c_state.Context["cxt2"])
		fmt.Println("cxt3: ", string(c_state.Context["cxt3"].([]byte)))

		fmt.Println("------------")
	}
}
