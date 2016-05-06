package test

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"mzlib"
	"testing"
)

func TestMongoModel(t *testing.T) {
	fmt.Println("MongoModel test begin.......")
	mm := mzlib.NewMongoModel()
	mm.Uid = "778899"
	mm.Name = "xxxx"
	mm.Put()

	mm2 := mzlib.NewMongoModel()
	mm2.Get("778899")
	if mm.Name != mm2.Name {
		panic("assert fail!!!")
	}

	mm2.Delete()

	for i := 0; i < 10; i++ {
		mm3 := mzlib.NewMongoModel()
		mm3.Uid = fmt.Sprintf("000%d", i)
		mm3.Name = "papapa"
		mm3.Insert()
	}

	mm4 := mzlib.NewMongoModel()
	mm4.FindOne(bson.M{"uid": "0000"})
	fmt.Printf("mm4=============%#v\n", mm4)

	mm5 := mzlib.NewMongoModel()
	list := mm5.Find(bson.M{"name": "papapa"})
	fmt.Printf("list=============%#v\n", list)

	for i, v := range list {
		vv, ok := v.(*mzlib.MongoModel)
		fmt.Println("OK============", ok)
		fmt.Printf("list %d=============%#v====%v===%s\n", i, v, ok, vv.Name)
	}

	for i := 0; i < 10; i++ {
		mmdel := mzlib.NewMongoModel()
		uid := fmt.Sprintf("000%d", i)
		mmdel.Get(uid)
		mmdel.Delete()
	}

	fmt.Println("MongoModel test end.........")
}
