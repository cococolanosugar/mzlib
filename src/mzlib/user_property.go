package mzlib

import (
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
)

type UserProperty struct {
	Pk    string `json:"-" bson:"-"`
	Uid   string `json:"uid" bson:"uid"`
	Gold  int    `json:"gold" bson:"gold"`
	Coin  int    `json:"coin" bson:"coin"`
	Name  string `json:"name" bson:"name"`
	Title string `json:"Title" bson:"-"`
}

func NewUserProperty() *UserProperty {
	up := new(UserProperty)
	up.Pk = "Uid"
	return up
}

func UserPropertyGet(uid string) *UserProperty {
	up := NewUserProperty()
	res := app.Get(up, uid)
	res2, ok := res.(*UserProperty)
	if ok {
		return res2
	}
	return nil
}

func (up *UserProperty) Put() {
	app.Put(up)
}

func (up *UserProperty) DoPut() {
	app.DoPut(up)
}

func (up *UserProperty) ToJson() string {
	js, _ := json.Marshal(up)
	return string(js)
}

func (up *UserProperty) ToBson() string {
	bs, _ := bson.Marshal(up)
	return string(bs)
}

func (up *UserProperty) Load(data interface{}) {
	data_bson, ok := data.(bson.M)
	if !ok {
		panic("load data invalid!!!")
	}
	b, _ := json.Marshal(data_bson)
	json.Unmarshal(b, up)
}
