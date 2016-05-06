package mzlib

import (
	"encoding/json"
	"fmt"
	"mzlib/client"
	"reflect"
	"strings"
)

const (
	USE_PIER = true
)

type App struct {
	redis_storage *StorageRedis
	mongo_storage *StorageMongo
	top_redis     *client.MZRedis
	log_mongo     *StorageMongo
	pier          *Pier

	RedisStorage *StorageRedis
	MongoStorage *StorageMongo
	TopRedis     *client.MZRedis
	LogMongo     *StorageMongo
	Pier         *Pier
}

func NewApp() *App {
	app := new(App)
	app.Init()
	return app
}

func (app *App) Init() {
	InitConfig()

	app.redis_storage = NewStorageRedis(redis_conf)
	app.mongo_storage = NewStorageMongo(mongo_conf)
	app.log_mongo = NewStorageMongo(mongo_log_conf)
	app.top_redis = client.NewMZRedisClient(redis_top_conf.Host, redis_top_conf.Port, redis_top_conf.Index, redis_top_conf.Password)
	app.pier = NewPier()

	app.RedisStorage = app.redis_storage
	app.MongoStorage = app.mongo_storage
	app.LogMongo = app.log_mongo
	app.TopRedis = app.top_redis
	app.Pier = app.pier
}

func (app *App) Get(model interface{}, pk_value string) interface{} {
	//从pier中获取
	pier_res := app.pier.Get(model, pk_value)
	fmt.Println("AppGet    pier_res=====", pier_res, reflect.TypeOf(pier_res))
	if pier_res != nil {
		return pier_res
	}
	//从redis中获取
	redis_res := app.redis_storage.Get(model, pk_value)
	fmt.Println("AppGet    redis_res=====", redis_res, reflect.TypeOf(redis_res))
	if redis_res != "" {
		js, ok := redis_res.(string)
		if ok {
			json.Unmarshal([]byte(js), model)
			return model
		}
		return nil
	}
	//从mongo中获取
	mongo_res := app.mongo_storage.Get(model, pk_value)
	fmt.Println("AppGet    mongo_res=====", mongo_res, reflect.TypeOf(mongo_res))
	if mongo_res != nil {
		//mongo_res 返回的是bson.M 类型,之所以不用 bson.Marshal bson.Unmarshal 是因为初始化的时候给Pk赋的值会被置为nil然后再去DoPut就会找不到Pk
		js, _ := json.Marshal(mongo_res)
		json.Unmarshal(js, model)
		return model
	}
	return nil
}

func (app *App) Put(model interface{}) {
	if app.pier.use {
		app.pier.Add(model)
	} else {
		app.redis_storage.Set(model)
	}
}

func (app App) DoPut(model interface{}) {
	app.redis_storage.Set(model)
	app.mongo_storage.Set(model)
}

type Pier struct {
	put_data map[string]map[string]interface{}
	store    StorageRedis
	get_data map[string]map[string]interface{}
	use      bool
	Use      bool
}

func NewPier() *Pier {
	pier := new(Pier)
	pier.use = USE_PIER
	pier.Use = USE_PIER
	pier.put_data = make(map[string]map[string]interface{})
	pier.get_data = make(map[string]map[string]interface{})
	return pier
}

func (pier *Pier) Add(user_model interface{}) {
	if !pier.use {
		return
	}

	//uid := user_model.Uid
	uid := reflect.ValueOf(user_model).Elem().FieldByName("Uid").String()
	if uid == "" {
		panic("invalid uid!!!")
	}

	//add to put_data
	fmt.Println(1111)
	if _, ok := pier.put_data[uid]; !ok {
		pier.put_data[uid] = map[string]interface{}{}
		fmt.Println(11110000)
	}
	fmt.Println(2222)
	key_put := pier.get_class_key(user_model)
	if _, ok := pier.put_data[uid][key_put]; !ok {
		pier.put_data[uid][key_put] = user_model
		fmt.Println(22220000)
	}
	fmt.Println(3333)
	//add to get_data
	if _, ok := pier.get_data[uid]; !ok {
		pier.get_data[uid] = map[string]interface{}{}
		fmt.Println(33330000)
	}
	fmt.Println(4444)
	key_get := pier.get_class_key(user_model)
	if _, ok := pier.get_data[uid][key_get]; !ok {
		pier.get_data[uid][key_get] = user_model
		fmt.Println(44440000)
	}
	fmt.Println(5555)
	fmt.Println("Pier Add put_data=======", pier.put_data)
	fmt.Println("Pier Add get_data=======", pier.get_data)
}

func (pier *Pier) AddGetData(user_model interface{}) {
	if !pier.use {
		return
	}

	//add to get_data
	//uid := user_model.Uid
	uid := reflect.ValueOf(user_model).Elem().FieldByName("Uid").String()
	if uid == "" {
		panic("invalid uid!!!")
	}

	if _, ok := pier.get_data[uid]; !ok {
		pier.get_data[uid] = map[string]interface{}{}
	}
	key_get := pier.get_class_key(user_model)
	if _, ok := pier.get_data[uid][key_get]; !ok {
		pier.get_data[uid][key_get] = user_model
	}
}

func (pier *Pier) Save() {
	if !pier.use {
		return
	}

	if len(pier.put_data) == 0 {
		return
	}
	for k, vm := range pier.put_data {
		var um_list []interface{}
		//for k1, v1 := range pier.put_data[k] {
		for _, v1 := range vm {
			um_list = append(um_list, v1)
		}
		pier.store.MsetUserModel(k, um_list)
	}
	pier.Clear()
}

func (pier *Pier) Clear() {
	if !pier.use {
		return
	}

	if len(pier.put_data) == 0 {
		return
	}
	for k, _ := range pier.put_data {
		delete(pier.put_data, k)
	}
}

func (pier *Pier) Get(user_model interface{}, uid string) interface{} {
	if !pier.use {
		return nil
	}

	um_map, ok := pier.get_data[uid]
	fmt.Println(1111)
	if !ok {
		fmt.Println(11110000)
		return nil
	}
	fmt.Println(2222)
	if len(um_map) == 0 {
		fmt.Println(22220000)
		return nil
	}
	fmt.Println(3333)
	for _, v := range um_map {
		if reflect.TypeOf(v) == reflect.TypeOf(user_model) {
			fmt.Println(33330000)
			fmt.Println("Pier Get put_data=======", pier.put_data)
			fmt.Println("Pier Get get_data=======", pier.get_data)
			return v
		}
	}
	fmt.Println(4444)
	fmt.Println("Pier Get put_data=======", pier.put_data)
	fmt.Println("Pier Get get_data=======", pier.get_data)
	return nil
}

func (pier *Pier) get_class_key(user_model interface{}) string {
	class := reflect.TypeOf(user_model).Elem()
	class_name := class.Name()
	//class_name = strings.ToLower(class_name)
	uid := reflect.ValueOf(user_model).Elem().FieldByName("Uid").String()
	if uid == "" {
		panic("invalid uid!!!")
	}
	key := strings.Join([]string{class_name, uid}, ":")
	key = strings.ToLower(key)
	return key
}

var app = NewApp()
var AppIns = app
