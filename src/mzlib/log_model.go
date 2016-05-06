package mzlib

import (
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	//"gopkg.in/redis.v3"
	"reflect"
	"time"
)

// const (
// 	TIME_FORMAT = "2006-01-02 15:04:05"
// )

//logmodel 存储log数据 写一次后基本不修改 有查询需求
type LogModel struct {
	Pk       string `json:"-" bson:"-"`
	Uid      string `json:"uid" bson:"uid`
	Name     string `json:"name" bson:"name"`
	DateTime string `json:"datetime" bson:"datetime"`
}

/*
	Use Case:
	type LoginLog struct {
		mzlib.LogModel `bson:",inline"`
		Gold           int `json:"gold" bson:"gold"`
	}
	log := new(LoginLog)
	log.Pk = "Uid"
	log.Uid = "002"
	log.Name = "xushicai"
	log.DateTime = time.Now().Format("2006-01-02 15:04:05")
	log.Gold = 100
	LogModelSet(log)

	mongo>>>>db.loginlog.find({})
*/
//往logdb中插入一条log记录
func LogModelSet(log interface{}) {
	//启动一个gorutine插入日志
	go func() {
		app.log_mongo.Insert(log)
	}()
}

func NewLogModel() *LogModel {
	lm := new(LogModel)
	lm.Pk = "Uid"
	lm.DateTime = time.Now().Format(TIME_FORMAT)
	return lm
}

func (lm *LogModel) Create() {
	//TO DO
}

func (lm *LogModel) Put() {
	app.log_mongo.Insert(lm)
}

func (lm *LogModel) Find(query bson.M) []interface{} {
	list := app.log_mongo.Find(lm, query)
	fmt.Printf("LogModel  Find  list==========%#v", list)
	return list
}

func (lm *LogModel) Aggregate(statement bson.M) {
	//TO DO
	//app.log_mongo.Aggregate(lm, statement)
}

func (lm *LogModel) Load(data interface{}) {
	fmt.Println("LogModel Load data=======", data, reflect.TypeOf(data))
	data_bson, ok := data.(bson.M)
	if !ok {
		panic("load data invalid!!!")
	}
	b, _ := json.Marshal(data_bson)
	fmt.Println("LogModel Load b=======", b, reflect.TypeOf(b))
	fmt.Printf("LogModel Load b=======%s\n", b)
	json.Unmarshal(b, lm)
	fmt.Printf("LogModel Load bm=======%#v\n", lm)
}

func (lm *LogModel) Update() {
	//TO DO
}

func (lm *LogModel) ToJson() {
	json.Marshal(lm)
}

func (lm *LogModel) ToBson() interface{} {
	v, _ := bson.Marshal(lm)
	fmt.Printf("LogModel ToBson v==============%s\n", v)
	return v
}
