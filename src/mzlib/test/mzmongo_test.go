package test

import (
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"mzlib/client"
	"reflect"
	"testing"
)

func TestMZMongoDB(t *testing.T) {
	return

	fmt.Println("mongo test begin.......")

	mongo := client.NewMZMongoClient("mongodb://127.0.0.1:27017")
	type Person struct {
		Name  string
		Phone string
	}
	type Men struct {
		Persons []Person
	}

	query := func(sess *mgo.Session) error {
		fmt.Println("i am in query")
		collection := sess.DB("test").C("people")
		//集合数量
		cnt, err := collection.Count()
		checkErr(err)
		fmt.Println("collection objects count:", cnt)
		//单条查询
		res := Person{}
		query_json := bson.M{"phone": "15618973660"}
		err = collection.Find(query_json).One(&res)
		checkErr(err)
		s, _ := json.Marshal(res)
		fmt.Printf("single res is :%s\n", s)
		//多条查询
		var personAll Men
		query_json = bson.M{"phone": "15618973660"}
		err = collection.Find(query_json).All(&personAll.Persons)
		s, _ = json.Marshal(personAll)
		fmt.Printf("personall is :%s", s)

		//迭代查询
		//var men Men
		iter := collection.Find(nil).Iter()
		for iter.Next(&res) {
			fmt.Println("name === ", res.Name)
		}

		return nil
	}

	insert := func(sess *mgo.Session) error {
		p := &Person{
			Name:  "papapa",
			Phone: "15618973660",
		}

		p1 := &Person{
			Name:  "kakaka",
			Phone: "156189736600",
		}

		p2 := &Person{
			Name:  "dadada",
			Phone: "156189736600",
		}

		p3 := &Person{
			Name:  "ooooo",
			Phone: "156189736600",
		}

		p4 := &Person{
			Name:  "xxxx",
			Phone: "156189736600",
		}

		collection := sess.DB("test").C("people")
		err := collection.Insert(p, p1, p2, p3, p4, &Person{"hehehe", "18516604069"})
		checkErr(err)
		return err
	}

	update := func(sess *mgo.Session) error {
		collection := sess.DB("test").C("people")
		query_json := bson.M{"name": "kakaka"}
		modify_json := bson.M{"$set": bson.M{"name": "kakaka999"}}
		collection.Update(query_json, modify_json)

		query_json = bson.M{"name": "dadada"}
		modify_json = bson.M{"$set": bson.M{"phone": "00000000000"}}
		collection.Update(query_json, modify_json)

		query_json = bson.M{"name": "ooooo"}
		modify_json = bson.M{"phone": "9999999", "name": "uuuuuuuuu"}
		collection.Update(query_json, modify_json)

		return nil
	}

	remove := func(sess *mgo.Session) error {
		collection := sess.DB("test").C("people")
		n, err := collection.RemoveAll(bson.M{"name": "xxxx"})
		fmt.Println("delete cnt :", n)
		if false {
			collection.RemoveAll(nil)
		}
		return err
	}

	mongo.Execute(insert)
	mongo.Execute(query)
	mongo.Execute(insert)
	mongo.Execute(update)
	mongo.Execute(remove)
	fmt.Println("mongo test end.........")
}

func TestMZMongoDBV2(t *testing.T) {
	fmt.Println("mongo test begin.......")

	mongo := client.NewMZMongoClient("mongodb://127.0.0.1:27017")
	collection_name := "people"
	pk := "name"
	pk_value := "kakaka999"

	// test insert
	for i := 0; i < 10; i++ {
		mongo.Insert(collection_name, bson.M{pk: pk_value, "phone": "15618973660"})
		mongo.Insert(collection_name, bson.M{"name": "papapa", "phone": "15618973660"})
	}

	//test get
	res := mongo.Get(collection_name, pk, pk_value)
	fmt.Printf("res=======%s\n", res)

	res1 := mongo.Get(collection_name, "_id", "56f50da664ad243a18cdef4a")
	fmt.Printf("res1=======%s\n", res1)

	//test find
	res2 := mongo.Find(collection_name, bson.M{"name": "papapa"})
	fmt.Printf("res2=======%s\n", res2)
	fmt.Printf("res2 type=======%s\n", reflect.TypeOf(res2))
	for i, v := range res2 {
		tmp, _ := json.Marshal(v)
		fmt.Printf("result2===%d====%s\n", i, tmp)
	}

	//test delete
	mongo.Delete(collection_name, "name", "papapa")
	res3 := mongo.Find(collection_name, bson.M{"name": "papapa"})
	Assert(len(res3), len(res2)-1, t)

	mongo.DeleteAll(collection_name, bson.M{"name": "papapa"})
	res4 := mongo.Find(collection_name, bson.M{"name": "papapa"})
	Assert(len(res4), 0, t)

	mongo.DeleteAll(collection_name, bson.M{"name": "cccc"})
	mongo.Insert(collection_name, bson.M{"name": "cccc", "phone": "xxxxxx"})
	mongo.Update(collection_name, "name", bson.M{"name": "cccc", "phone": "gggggg"})
	res5 := mongo.Get(collection_name, "name", "cccc")
	res6, ok := res5.(bson.M)
	fmt.Println("ok======", ok, reflect.TypeOf(res5).Kind())
	Assert(res6["phone"], "gggggg", t)
	//fmt.Printf("res5=======%s\n", res5)

	mongo.Upsert(collection_name, "name", bson.M{"name": "cccc", "phone": "hhhhhh"})
	res7 := mongo.Get(collection_name, "name", "cccc")
	res8, _ := res7.(bson.M)
	Assert(res8["phone"], "hhhhhh", t)

	mongo.Upsert(collection_name, "name", bson.M{"name": "kkkkkkk", "phone": "hhhhhh"})
	res9 := mongo.Get(collection_name, "name", "kkkkkkk")
	res10, _ := res9.(bson.M)
	Assert(res10["phone"], "hhhhhh", t)

	mongo.UpdateAll(collection_name, "name", bson.M{"name": "dadada", "phone": "0000001111"})
	res11 := mongo.Find(collection_name, bson.M{"name": "dadada"})
	for i, v := range res11 {
		tmp, _ := json.Marshal(v)
		fmt.Printf("result===%d====%s\n", i, tmp)
	}

	idxs := mongo.GetIndexes(collection_name)
	for k, v := range idxs {
		s, _ := json.Marshal(v)
		fmt.Printf("idxs========%d=====%s\n", k, s)
	}
	mongo.EnsureIndex(collection_name, "name")
	mongo.EnsureIndex(collection_name, "phone")
	idxs1 := mongo.GetIndexes(collection_name)
	for k, v := range idxs1 {
		s, _ := json.Marshal(v)
		fmt.Printf("idxs========%d=====%s\n", k, s)
	}

	mongo.DropIndex(collection_name, "name")
	idxs2 := mongo.GetIndexes(collection_name)
	for k, v := range idxs2 {
		s, _ := json.Marshal(v)
		fmt.Printf("idxs========%d=====%s\n", k, s)
	}

	fmt.Println("mongo test end.........")
}
