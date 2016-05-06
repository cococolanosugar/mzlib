package mzlib

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"reflect"
)

type UserName struct {
	Pk   string `json:"-" bson:"-"`
	Name string `json:"name" bson:"name"`
	Uid  string `json:"uid" bson:"uid"`
}

func NewUserName() *UserName {
	un := new(UserName)
	un.Pk = "Uid"
	return un
}

func UserNameGetByName(name string) *UserName {
	if name == "" {
		return nil
	}

	query := bson.M{}
	if name != "" {
		query["name"] = name
	}

	un := NewUserName()
	mongo_res := app.mongo_storage.FindOne(un, query)
	if mongo_res != nil {
		un.Load(mongo_res)
		return un
	}
	return nil
}

func UserNameGetByUid(uid string) *UserName {
	if uid == "" {
		return nil
	}

	query := bson.M{}
	if uid != "" {
		query["uid"] = uid
	}

	un := NewUserName()
	mongo_res := app.mongo_storage.FindOne(un, query)
	if mongo_res != nil {
		un.Load(mongo_res)
		return un
	}
	return nil
}

func UserNameSet(name string, uid string) error {
	if name == "" || uid == "" {
		return errors.New("name or uid is null!")
	}
	un := NewUserName()
	un.Uid = uid
	un.Name = name
	app.mongo_storage.Insert(un)
	return nil
}

func (un *UserName) Load(data interface{}) {
	fmt.Println("UserName Load data=======", data, reflect.TypeOf(data))
	data_bson, ok := data.(bson.M)
	if !ok {
		panic("load data invalid!!!")
	}
	b, _ := json.Marshal(data_bson)
	fmt.Println("UserName Load b=======", b, reflect.TypeOf(b))
	fmt.Printf("UserName Load b=======%s\n", b)
	json.Unmarshal(b, un)
	fmt.Printf("UserName Load bm=======%#v\n", un)
}
