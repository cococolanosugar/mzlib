package mzlib

import (
	//"encoding/json"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"hash/crc32"
	"mzlib/client"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	R2M_SET    = "r2m_set"
	R2M        = true
	EXPIRATION = time.Hour * 24
)

//storage
type Storage interface {
	Get(id string) interface{}
	Set(obj interface{}) error
}

//redis storage
type StorageRedis struct {
	redis_list []*client.MZRedis
}

func NewStorageRedis(redis_conf RedisConfig) *StorageRedis {
	rs := new(StorageRedis)
	rs.Init(redis_conf)
	return rs
}

func (rs *StorageRedis) Init(redis_conf RedisConfig) {
	redis_client := client.NewMZRedisClient(redis_conf.Host, redis_conf.Port, redis_conf.Index, redis_conf.Password)
	rs.redis_list = append(rs.redis_list, redis_client)
}

func (rs *StorageRedis) Get(model_class interface{}, pk string) interface{} {
	key := rs.get_class_key(model_class, pk)
	value := rs.redis_list[rs.hash(key)%uint32(len(rs.redis_list))].Get(key)
	fmt.Println("StorageRedis Get key =======", key)
	fmt.Println("StorageRedis Get value =======", value)
	//json.Unmarshal([]byte(value), model_class)
	return value
}

func (rs *StorageRedis) Set(obj interface{}) {
	key := rs.get_obj_key(obj)

	/*暂时换成下面方式 方便继承
	method := reflect.ValueOf(obj).MethodByName("ToJson")

	value := method.Call([]reflect.Value{})
	final_value := value[0].Interface()
	*/
	final_value, err := json.Marshal(obj)
	redis_client := rs.redis_list[rs.hash(key)%uint32(len(rs.redis_list))]
	if R2M {
		redis_client.Sadd(R2M_SET, key)
	}
	fmt.Println("StorageRedis Set key=========", key, err)
	//fmt.Println("StorageRedis Set value=======", value, len(value), value[0].Interface())
	//fmt.Printf("StorageRedis Set value[0]=========%s\n", value[0].Interface())
	fmt.Printf("StorageRedis Set final_value=========%s\n", final_value)
	redis_client.Set(key, final_value, EXPIRATION)
}

func (rs *StorageRedis) GetUseModel(uid string) {

}

func (rs *StorageRedis) SetUserModel(obj interface{}) {

}

func (rs *StorageRedis) MsetUserModel(uid string, user_model_set []interface{}) {

}

func (rs *StorageRedis) Delete(obj interface{}) {
	key := rs.get_obj_key(obj)
	redis_client := rs.redis_list[rs.hash(key)%uint32(len(rs.redis_list))]
	redis_client.Srem(R2M_SET, key)
	redis_client.Delete(key)
}

func (rs *StorageRedis) DeleteUserModel(obj interface{}) {
	key := rs.get_obj_key(obj)
	redis_client := rs.redis_list[rs.hash(key)%uint32(len(rs.redis_list))]
	redis_client.Srem(R2M_SET, key)
	redis_client.Delete(key)
}

func (rs *StorageRedis) get_class_key(model_class interface{}, pk_value string) string {
	class := reflect.TypeOf(model_class).Elem()
	class_name := class.Name()
	key := strings.Join([]string{class_name, pk_value}, ":")
	key = strings.ToLower(key)
	return key
}

func (rs *StorageRedis) get_obj_key(model_class interface{}) string {
	class := reflect.TypeOf(model_class).Elem()
	pk := reflect.ValueOf(model_class).Elem().FieldByName("Pk").String()
	if pk == "" {
		panic("invalid pk")
	}
	uid := reflect.ValueOf(model_class).Elem().FieldByName(pk).String()
	if uid == "" {
		panic("invalid pk")
	}
	class_name := class.Name()
	key := strings.Join([]string{class_name, uid}, ":")
	key = strings.ToLower(key)
	fmt.Println("StorageRedis get_obj_key class===========", class)
	fmt.Println("StorageRedis get_obj_key class_name===========", class_name)
	fmt.Println("StorageRedis get_obj_key pk===========", pk)
	fmt.Println("StorageRedis get_obj_key key===========", key)
	return key
}

func (rs *StorageRedis) hash(key string) uint32 {
	if len(rs.redis_list) == 1 {
		return 1
	}
	return crc32.ChecksumIEEE([]byte(key))
}

//mongo storage
type StorageMongo struct {
	mongo *client.MZMongo
}

func NewStorageMongo(mongo_conf MongodbConfig) *StorageMongo {
	rm := new(StorageMongo)
	rm.Init(mongo_conf)
	return rm
}

func (ms *StorageMongo) Init(mongo_conf MongodbConfig) {
	//mongodb://127.0.0.1:27017
	mongo_url := "mongodb://"
	if mongo_conf.Uname != "" && mongo_conf.Pass != "" {
		mongo_url += mongo_conf.Uname + ":" + mongo_conf.Pass + "@"
	}
	mongo_url += mongo_conf.Host + ":" + strconv.Itoa(mongo_conf.Port) + "/" + mongo_conf.Dbname
	ms.mongo = client.NewMZMongoClient(mongo_url)
}

func (ms *StorageMongo) Get(obj interface{}, pk_value string) interface{} {
	collection_name := ms.get_collection_name(obj)
	pk := ms.get_pk(obj)
	data := ms.mongo.Get(collection_name, pk, pk_value)
	fmt.Printf("StorageMongo Get data======%#v\n", data)
	return data
	//res := reflect.ValueOf(obj).MethodByName("Load").Call([]reflect.Value{reflect.ValueOf(data)})
	//return res
}

func (ms *StorageMongo) Set(obj interface{}) {
	collection_name := ms.get_collection_name(obj)
	pk := ms.get_pk(obj)
	/*暂时换成下面方式 方便继承
	value := reflect.ValueOf(obj).MethodByName("ToBson").Call([]reflect.Value{})
	final_value := value[0].Interface()
	*/

	//新方式
	final_value, ok := bson.Marshal(obj)

	fmt.Println("StorageMongo Set pk======", pk)
	fmt.Printf("StorageMongo Set final_value======%s=====%v\n", final_value, ok)
	bs := bson.M{}
	bson.Unmarshal(final_value, bs)

	fmt.Println("StorageMongo Set bs======", bs)
	fmt.Println("StorageMongo Set collection_name=====pk=====", collection_name, pk)
	ms.mongo.Upsert(collection_name, pk, bs)
	fmt.Println("StorageMongo Set collection_name=====pk=====", collection_name, pk)
}

func (ms *StorageMongo) Insert(obj interface{}) {
	collection_name := ms.get_collection_name(obj)
	//ms.mongo.Insert(collection_name, obj.ToJson())

	/*暂时换成下面方式 方便继承
	value := reflect.ValueOf(obj).MethodByName("ToBson").Call([]reflect.Value{})
	final_value := value[0].Interface()
	*/

	//新方式
	final_value, err := bson.Marshal(obj)
	fmt.Println("StorageMongo Insert final_value======", final_value, err)

	bs := bson.M{}
	bson.Unmarshal(final_value, bs)
	ms.mongo.Insert(collection_name, bs)
	fmt.Println("StorageMongo Insert collection_name==========", collection_name)
}

func (ms *StorageMongo) Find(obj interface{}, query bson.M) []interface{} {
	collection_name := ms.get_collection_name(obj)
	res_list := ms.mongo.Find(collection_name, query)
	var final_list []interface{}
	if len(res_list) > 0 {
		for _, v := range res_list {
			o := reflect.New(reflect.TypeOf(obj).Elem()).Interface()
			fmt.Println("**************", reflect.TypeOf(obj), reflect.TypeOf(obj).Elem(), reflect.TypeOf(o))
			value := reflect.ValueOf(o).MethodByName("Load").Call([]reflect.Value{reflect.ValueOf(v)})
			fmt.Println("StorageMongo Find value======", value)
			final_list = append(final_list, o)
		}
	}
	return final_list
}

func (ms *StorageMongo) FindOne(obj interface{}, query bson.M) interface{} {
	collection_name := ms.get_collection_name(obj)
	v := ms.mongo.FindOne(collection_name, query)
	//value := reflect.ValueOf(obj).MethodByName("Load").Call([]reflect.Value{reflect.ValueOf(v)})
	fmt.Printf("StorageMongo FindOne res========%#v\n", v)
	//fmt.Println("StorageMongo FindOne value======", value)
	return v
}

func (ms *StorageMongo) FindAndModify(obj interface{}, query bson.M, update bson.M) interface{} {
	collection_name := ms.get_collection_name(obj)
	v := ms.mongo.FindAndModify(collection_name, query, update)
	fmt.Printf("StorageMongo FindAndModify res========%#v\n", v)
	return v
}

func (ms *StorageMongo) Delete(obj interface{}) {
	collection_name := ms.get_collection_name(obj)
	pk := ms.get_pk(obj)
	pk_value := ms.get_pk_value(obj)
	ms.mongo.Delete(collection_name, pk, pk_value)
	fmt.Println("StorageMongo Delete pk,pk_value=======", pk, pk_value)
}

func (ms *StorageMongo) OrderBy() {

}

func (ms *StorageMongo) get_collection_name(obj interface{}) string {
	class_name := reflect.TypeOf(obj).Elem().Name()
	collection_name := strings.ToLower(class_name)
	return collection_name
}

func (ms *StorageMongo) get_pk(obj interface{}) string {
	pk := reflect.ValueOf(obj).Elem().FieldByName("Pk").String()
	if pk == "" {
		panic("invalid pk")
	}
	pk = strings.ToLower(pk)
	return pk
}

func (ms *StorageMongo) get_pk_value(obj interface{}) string {
	pk := reflect.ValueOf(obj).Elem().FieldByName("Pk").String()
	if pk == "" {
		panic("invalid pk")
	}
	pk_value := reflect.ValueOf(obj).Elem().FieldByName(pk).String()
	return pk_value
}
