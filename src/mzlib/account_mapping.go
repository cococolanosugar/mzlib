package mzlib

import (
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"time"
)

type AccountMapping struct {
	Pk          string `json:"-" bson:"-"`
	Pid         string `json:"pid" bson:"pid"`
	Openid      string `json:"openid" bson:"openid"`
	Platform    string `json:"platform" bson:"platform"`
	Uid         string `json:"uid" bson:"uid"`
	AccessToken string `json:"accesstoken" bson:"accesstoken"`
	CreateAt    string `json:"createat" bson:"createat"`
}

func NewAccountMapping(pid string) *AccountMapping {
	am := new(AccountMapping)
	am.Pk = "Pid"
	am.Pid = pid
	am.Uid = NextUid()
	am.CreateAt = time.Now().Format(TIME_FORMAT)
	return am
}

func CreateAccountMapping(pid string, uid string) *AccountMapping {
	am := new(AccountMapping)
	am.Pk = "Pid"
	am.Pid = pid
	am.Uid = uid
	am.CreateAt = time.Now().Format(TIME_FORMAT)
	return am
}

func AccountMappingGet(pid string) *AccountMapping {
	am := NewAccountMapping(pid)
	//从pier中获取
	pier_res := app.pier.Get(am, pid)
	fmt.Println("AccountMappingGet    pier_res=====", pier_res, reflect.TypeOf(pier_res))
	if pier_res != nil {
		am, ok := pier_res.(*AccountMapping)
		if ok {
			return am
		}
		return nil
	}
	//从redis中获取
	redis_res := app.redis_storage.Get(am, pid)
	fmt.Println("AccountMappingGet    redis_res=====", redis_res, reflect.TypeOf(redis_res))
	if redis_res != "" {
		js, ok := redis_res.(string)
		if ok {
			json.Unmarshal([]byte(js), am)
			return am
		}
		return nil
	}
	//从mongo中获取
	mongo_res := app.mongo_storage.Get(am, pid)
	fmt.Println("AccountMappingGet    mongo_res=====", mongo_res, reflect.TypeOf(mongo_res))
	if mongo_res != nil {
		am.Load(mongo_res)
		return am
	}
	return nil
}

func (am *AccountMapping) Put() {
	if app.pier.use {
		app.pier.Add(am)
	} else {
		app.redis_storage.Set(am)
	}
}

func (am *AccountMapping) DoPut() {
	app.redis_storage.Set(am)
	app.mongo_storage.Set(am)
}

func (am *AccountMapping) Load(data interface{}) {
	fmt.Println("AccountMapping Load data=======", data, reflect.TypeOf(data))
	//data_json, ok := data.([]byte)
	data_bson, ok := data.(bson.M)
	if !ok {
		panic("load data invalid!!!")
	}
	b, _ := json.Marshal(data_bson)
	//json.Unmarshal([]byte(data_json), bm)
	fmt.Println("AccountMapping Load b=======", b, reflect.TypeOf(b))
	fmt.Printf("AccountMapping Load b=======%s\n", b)
	json.Unmarshal(b, am)
	fmt.Printf("AccountMapping Load bm=======%#v\n", am)
}
