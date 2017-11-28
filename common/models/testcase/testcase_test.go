package testcase

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

func TestObjectId(t *testing.T) {
	a := bson.NewObjectId()
	fmt.Printf("new oid: %v\n", a)
	fmt.Printf("str: %s\n", a)
	fmt.Printf("hex: %s\n", a.Hex())

	b := a.Hex()
	c := bson.ObjectIdHex(b)
	fmt.Println(a == c)

	var d bson.ObjectId
	// comparable with blank string
	fmt.Println(d == "")
}
