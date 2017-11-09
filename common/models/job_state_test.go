package models

import (
	"fmt"
	"testing"
)

func TestJobState(t *testing.T) {

	var jobStates = []JobState{
		{
			ActionRunId: "action_run_1",
			State:       "running",
			Results: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			Logs: []byte(`this is logs
                this is logs
                thi is logs`),
		},
	}

	for _, state := range jobStates {
		// test marshalling
		data, err := state.ToJson()
		if err != nil {
			t.Errorf("failed to marshal: %v\n", err)
		}
		fmt.Println(string(data))

		// test unmarshalling
		c_state, err := JobStateFromJson(data)
		if err != nil {
			t.Errorf("cannot unmarshal: %v\n", err)
		}

		fmt.Println("result item: ", c_state.Results["key2"])
		fmt.Println("logs: ", string(c_state.Logs))
		fmt.Println("-----------")
	}
}
