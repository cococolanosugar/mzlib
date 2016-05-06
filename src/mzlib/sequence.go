package mzlib

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"reflect"
)

const (
	SEQUENCE_COLLECTION = "sequence"
)

type Sequence struct {
	ID  string `json:"_id" bson:"_id"`
	Seq int64  `json:"seq" bson:"seq"`
}

func NextUid() string {
	sequence := new(Sequence)
	query := bson.M{"_id": "userid"}
	update := bson.M{"$inc": bson.M{"seq": 1}}
	v := app.mongo_storage.FindAndModify(sequence, query, update)
	fmt.Println("v=============", v, reflect.TypeOf(v))
	// bs, ok := v.(bson.M)
	// if !ok {
	// 	panic("invalid type of value")
	// }
	bsbyte, _ := bson.Marshal(v)
	bson.Unmarshal(bsbyte, sequence)

	return fmt.Sprint(sequence.Seq)
}
