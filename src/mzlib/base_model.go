package mzlib

import (
	//"github.com/bitly/go-simplejson"
	//"strings"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	//"gopkg.in/redis.v3"
	"reflect"
	//"time"
)

//model interface
type Model interface {
	Get(pk string) interface{}
	Put(obj interface{})
	DoPut(obj interface{})
	Create() interface{}
	Delete(obj interface{})
	Load(data interface{})
	ToJson() interface{}
	ToBson() interface{}
}

//basemodle
type BaseModel struct {
	Pk   string                 `json:"-" bson:"-"`
	Uid  string                 `json:"uid" bson:"uid"`
	Info map[string]interface{} `json:"info" bson:"info"`
}

func NewBaseModel() *BaseModel {
	bm := new(BaseModel)
	bm.Pk = "Uid"
	return bm
}

func (bm *BaseModel) Get(pk_value string) interface{} {
	v := app.redis_storage.Get(bm, pk_value)
	// if !obj {
	// 	obj = app.mongo_storage.Get(bm, pk)
	// 	if obj {
	// 		obj.Put()
	// 	}
	// }
	fmt.Println("BaseModel Get v==========", v)
	return v
}

func (bm *BaseModel) Put() {
	app.redis_storage.Set(bm)
}

func (bm *BaseModel) DoPut() {
	app.redis_storage.Set(bm)
	app.mongo_storage.Set(bm)
}

func (bm *BaseModel) Create() interface{} {
	return nil
}

func (bm *BaseModel) Delete() {
	app.redis_storage.Delete(bm)
	app.mongo_storage.Delete(bm)
}

func (bm *BaseModel) Load(data interface{}) {
	return
	fmt.Println("BaseModel Load data=======", data, reflect.TypeOf(data))
	//data_json, ok := data.([]byte)
	data_bson, ok := data.(bson.M)
	if !ok {
		panic("load data invalid!!!")
	}
	b, _ := json.Marshal(data_bson)
	//json.Unmarshal([]byte(data_json), bm)
	fmt.Println("BaseModel Load b=======", b, reflect.TypeOf(b))
	fmt.Printf("BaseModel Load b=======%s\n", b)
	json.Unmarshal(b, bm)
	fmt.Printf("BaseModel Load bm=======%s\n", bm)
}

func (bm *BaseModel) ToJson() interface{} {
	v, _ := json.Marshal(bm)
	fmt.Printf("BaseModel ToJson v==============%s\n", v)
	return v
}

func (bm *BaseModel) ToBson() interface{} {
	v, _ := bson.Marshal(bm)
	fmt.Printf("BaseModel ToBson v==============%s\n", v)
	return v
}
