package job

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

func TestMarshal(t *testing.T) {
	j := &Job{
		RunId:    bson.NewObjectId(),
		ActionId: bson.NewObjectId(),
	}

	c, _ := j.Json()
	fmt.Println(string(c))

	j2, _ := FromJson(c)

	if j2.RunId != j.RunId {
		fmt.Println("converting back is failed")
		t.Fail()
	}

	fmt.Println(j2.RunId.Hex())
}

func TestCompareNil(t *testing.T) {
	j := &Job{
		RunId:    bson.NewObjectId(),
		ActionId: bson.NewObjectId(),
	}

	fmt.Printf("type of blank objectId is: %v\n", j.PrevActionId)

	fmt.Println(j.PrevActionId == "")
}
