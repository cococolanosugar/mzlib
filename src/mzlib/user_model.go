package mzlib

import (
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	//"gopkg.in/redis.v3"
	"reflect"
	//"time"
	"github.com/bitly/go-simplejson"
)

//usermodel
type UserModel struct {
	Pk         string                 `json:"-" bson:"-"`
	Pid        string                 `json:"pid" bson:"pid"`
	Uid        string                 `json:"uid" bson:"uid"`
	Openid     string                 `json:"openid" bson:"openid"`
	Name       string                 `json:"name" bson:"name"`
	Platform   string                 `json:"platform" bson:"platform"`
	Lv         int                    `json:"lv" bson:"lv"`
	Gold       int                    `json:"gold" bson:"gold"`
	Coin       int                    `json:"coin" bson:"coin"`
	Password   string                 `json:"password" bson:"password"`
	LashBattle map[string]interface{} `json:"lashbattle" bson:"lashbattle"`
	GameCtrl   *GameCtrl              `json:"gamectrl" bson:"gamectrl"`
}

func GetUserModelByPid(pid string) *UserModel {
	am := AccountMappingGet(pid)
	if am == nil {
		return nil
	}
	return UserModelGet(am.Uid)
}

func UserModelGet(uid string) *UserModel {
	um := NewUserModel()
	//从pier中获取
	pier_res := app.pier.Get(um, uid)
	fmt.Println("UserModelGet    pier_res=====", pier_res, reflect.TypeOf(pier_res))
	if pier_res != nil {
		um, ok := pier_res.(*UserModel)
		if ok {
			return um
		}
		return nil
	}
	//从redis中获取
	redis_res := app.redis_storage.Get(um, uid)
	fmt.Println("UserModelGet    redis_res=====", redis_res, reflect.TypeOf(redis_res))
	if redis_res != "" {
		js, ok := redis_res.(string)
		if ok {
			json.Unmarshal([]byte(js), um)
			return um
		}
		return nil
	}
	//从mongo中获取
	mongo_res := app.mongo_storage.Get(um, uid)
	fmt.Println("UserModelGet    mongo_res=====", mongo_res, reflect.TypeOf(mongo_res))
	if mongo_res != nil {
		um.Load(mongo_res)
		return um
	}
	return nil

}

func CreateUserModel(pid string, uid string, platform string, openid string) *UserModel {
	um := NewUserModel()
	um.Pid = pid
	um.Uid = uid
	um.Platform = platform
	um.Openid = openid
	//TODO Other init
	um.Gold = 18000
	return um
}

func NewUserModel() *UserModel {
	um := new(UserModel)
	um.Pk = "Uid"
	return um
}

func (um *UserModel) Get(uid string) interface{} {
	obj := app.pier.Get(um, uid)
	obj2 := app.redis_storage.Get(um, uid)
	obj3 := app.mongo_storage.Get(um, uid)
	fmt.Printf("obj============%#v\n", obj)
	fmt.Printf("obj2============%#v\n", obj2)
	fmt.Printf("obj3============%#v\n", obj3)
	return obj
}

func (um *UserModel) Put() {
	if app.pier.use {
		app.pier.Add(um)
		//app.redis_storage.Set(um)
	} else {
		app.redis_storage.Set(um)
	}
}

func (um *UserModel) DoPut() {
	app.redis_storage.Set(um)
	app.mongo_storage.Set(um)
}

func (um *UserModel) Delete() {
	app.redis_storage.Delete(um)
	app.mongo_storage.Delete(um)
}

func (um *UserModel) Load(data interface{}) {
	fmt.Println("UserModel Load data=======", data, reflect.TypeOf(data))
	//data_json, ok := data.([]byte)
	data_bson, ok := data.(bson.M)
	if !ok {
		panic("load data invalid!!!")
	}
	b, _ := json.Marshal(data_bson)
	//json.Unmarshal([]byte(data_json), bm)
	fmt.Println("UserModel Load b=======", b, reflect.TypeOf(b))
	fmt.Printf("UserModel Load b=======%s\n", b)
	json.Unmarshal(b, um)
	fmt.Printf("UserModel Load bm=======%#v\n", um)
}

func (um *UserModel) ToJson() interface{} {
	v, _ := json.Marshal(um)
	fmt.Printf("UserModel ToJson v==============%s\n", v)
	return v
}

func (um *UserModel) ToBson() interface{} {
	v, _ := bson.Marshal(um)
	fmt.Printf("UserModel ToBson v==============%s\n", v)
	return v
}

func (um *UserModel) WrapperInfo() *simplejson.Json {
	user_info := simplejson.New()
	user_info.Set("pid", um.Pid)
	user_info.Set("uid", um.Uid)
	user_info.Set("platform", um.Platform)
	user_info.Set("name", um.Name)
	user_info.Set("openid", um.Openid)
	user_info.Set("gold", um.Gold)
	user_info.Set("coin", um.Coin)
	user_info.Set("lv", um.Lv)
	return user_info
}
