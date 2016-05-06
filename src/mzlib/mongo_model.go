package mzlib

import (
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	//"gopkg.in/redis.v3"
	"reflect"
	//"time"
)

//mongomodel 只在mongo中存储 有查询需求
type MongoModel struct {
	Pk   string `json:"-" bson:"-"`
	Uid  string `json:"uid" bson:"uid`
	Name string `json:"name" bson:"name"`
}

func MongoModelGet(pk_value string) *MongoModel {
	mm := NewMongoModel()
	mongo_res := app.mongo_storage.Get(mm, pk_value)
	if mongo_res != nil {
		mm.Load(mongo_res)
		return mm
	}
	return nil
}

func NewMongoModel() *MongoModel {
	mm := new(MongoModel)
	mm.Pk = "Uid"
	return mm
}

func (mm *MongoModel) Get(pk_value string) {
	app.mongo_storage.Get(mm, pk_value)
}

func (mm *MongoModel) Put() {
	app.mongo_storage.Set(mm)
}

func (mm *MongoModel) Delete() {
	app.mongo_storage.Delete(mm)
}

func (mm *MongoModel) Find(query bson.M) []interface{} {
	list := app.mongo_storage.Find(mm, query)
	fmt.Printf("MongoModel Find list========%#v\n", list)
	return list
}

func (mm *MongoModel) Insert() {
	app.mongo_storage.Insert(mm)
}

func (mm *MongoModel) FindOne(query bson.M) {
	app.mongo_storage.FindOne(mm, query)
}

func (mm *MongoModel) FindAndModify() {
	//TO DO
	/*
		_, err := s.DB(db).C(collection).FindId(id).Apply(mgo.Change{
			Update:    bson.M{"$inc": bson.M{"seq": 1}},
			ReturnNew: true,
		}, &res)
	*/
}

func (mm *MongoModel) ToBson() interface{} {
	v, _ := bson.Marshal(mm)
	fmt.Printf("MongoModel ToBson v==============%s\n", v)
	return v
}

func (mm *MongoModel) Load(data interface{}) {
	fmt.Println("MongoModel Load data=======", data, reflect.TypeOf(data))
	data_bson, ok := data.(bson.M)
	if !ok {
		panic("load data invalid!!!")
	}
	b, _ := json.Marshal(data_bson)
	fmt.Println("MongoModel Load b=======", b, reflect.TypeOf(b))
	fmt.Printf("MongoModel Load b=======%s\n", b)
	json.Unmarshal(b, mm)
	fmt.Printf("MongoModel Load bm=======%#v\n", mm)
}
