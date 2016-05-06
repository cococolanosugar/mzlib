package mzlib

import (
	"encoding/json"
	"fmt"
	//"gopkg.in/mgo.v2/bson"
	//"gopkg.in/redis.v3"
	//"reflect"
	"time"
)

//tmpmodle临时存储数据
type TmpModel struct {
	Pk    string `json:"-"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

func TmpModelSet(key string, value string, ex time.Duration) *TmpModel {
	tm := NewTmpModel()
	tm.Key = key
	tm.Value = value
	tm.Put()
	return tm
}

func TmpModelGet(key string) string {
	tm := NewTmpModel()
	tm.Get(key)
	return tm.Value
}

func NewTmpModel() *TmpModel {
	tm := new(TmpModel)
	tm.Pk = "Key"
	return tm
}

func (tm *TmpModel) Get(pk string) {
	app.redis_storage.Get(tm, pk)
}

func (tm *TmpModel) Put() {
	fmt.Printf("TmpModel Put tm==========%#v\n", tm)
	app.redis_storage.Set(tm)
}

func (tm *TmpModel) Delete() {
	app.redis_storage.Delete(tm)
}

func (tm *TmpModel) ToJson() interface{} {
	v, _ := json.Marshal(tm)
	fmt.Printf("TmpModel ToJson v==============%s\n", v)
	return v
}
