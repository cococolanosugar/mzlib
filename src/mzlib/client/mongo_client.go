package client

import (
	//"fmt"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	//"os"
	"time"
)

const (
	DEFAULT_MGO_TIMEOUT = 300
	DEFAULT_CONCURRENT  = 128
	DEFAULT_MONGO_URL  = "mongodb://127.0.0.1:27017/test"
)

type MZMongo struct {
	dbname  string
	session *mgo.Session
	latch   chan *mgo.Session
}

func init() {

}

func NewMZMongoClient(mongom_url string) *MZMongo {
	mongo := new(MZMongo)
	mongo.Init(mongom_url)
	return mongo
}

func (m *MZMongo) Init(mongom_url string) {
	// create latch
	m.latch = make(chan *mgo.Session, DEFAULT_CONCURRENT)
	// connect db
	sess, err := mgo.Dial(mongom_url)
	if err != nil {
		log.Println("mongom: cannot connect to", mongom_url, err)
		//os.Exit(-1)
	}

	// set params
	sess.SetMode(mgo.Strong, true)
	sess.SetSocketTimeout(DEFAULT_MGO_TIMEOUT * time.Second)
	sess.SetCursorTimeout(0)
	m.session = sess

	// if m.dbname == "" {
	// 	m.dbname = "test"
	// }

	for k := 0; k < cap(m.latch); k++ {
		m.latch <- sess.Clone()
	}
}

func (m *MZMongo) Execute(f func(sess *mgo.Session) error) error {
	// latch control
	sess := <-m.latch
	defer func() {
		m.latch <- sess
	}()
	sess.Refresh()
	err := f(sess)
	if err != nil {
		log.Fatalln("mongodb err:", err)
	}
	return err
}

func (m *MZMongo) Get(collection_name, pk, pk_value string) interface{} {
	//根据文档的pk去读取文档
	var query_bson bson.M
	if pk == "_id" {
		query_bson = bson.M{pk: bson.ObjectIdHex(pk_value)}
	} else {
		query_bson = bson.M{pk: pk_value}
	}

	var res interface{}
	query_func := func(sess *mgo.Session) error {
		//单条查询
		collection := sess.DB(m.dbname).C(collection_name)
		err := collection.Find(query_bson).One(&res)
		if err != nil {
			return err
		}
		//s, _ := json.Marshal(res)
		//fmt.Printf("single res is :%s\n", s)
		//fmt.Println("=========", res)
		//res = s
		return nil
	}
	m.Execute(query_func)
	return res
}

func (m *MZMongo) Find(collection_name string, query_bson bson.M) []interface{} {
	var res []interface{}
	query_func := func(sess *mgo.Session) error {
		//多条查询
		collection := sess.DB(m.dbname).C(collection_name)
		err := collection.Find(query_bson).All(&res)
		if err != nil {
			return err
		}
		//s, _ := json.Marshal(res)
		//fmt.Printf("single res is :%s\n", s)
		//fmt.Println("=========", res)
		//res = s
		return nil
	}
	m.Execute(query_func)
	return res
}

func (m *MZMongo) FindOne(collection_name string, query_bson bson.M) interface{} {
	var res interface{}
	query_func := func(sess *mgo.Session) error {
		//多条查询
		collection := sess.DB(m.dbname).C(collection_name)
		err := collection.Find(query_bson).One(&res)
		if err != nil {
			return err
		}
		//s, _ := json.Marshal(res)
		//fmt.Printf("single res is :%s\n", s)
		//fmt.Println("=========", res)
		//res = s
		return nil
	}
	m.Execute(query_func)
	return res
}

func (m *MZMongo) Insert(collection_name string, data bson.M) {
	query_func := func(sess *mgo.Session) error {
		//单条插入
		collection := sess.DB(m.dbname).C(collection_name)
		err := collection.Insert(data)
		if err != nil {
			return err
		}
		return nil
	}
	m.Execute(query_func)
}

func (m *MZMongo) Delete(collection_name, pk, pk_value string) {
	var query_bson bson.M
	if pk == "_id" {
		query_bson = bson.M{pk: bson.ObjectIdHex(pk_value)}
	} else {
		query_bson = bson.M{pk: pk_value}
	}

	query_func := func(sess *mgo.Session) error {
		//单条删除
		collection := sess.DB(m.dbname).C(collection_name)
		//err := collection.Find(query_bson).All(&res)
		err := collection.Remove(query_bson)
		if err != nil {
			return err
		}
		//s, _ := json.Marshal(res)
		//fmt.Printf("single res is :%s\n", s)
		//fmt.Println("=========", res)
		//res = s
		return nil
	}
	m.Execute(query_func)
}

func (m *MZMongo) DeleteAll(collection_name string, query_bson bson.M) {

	query_func := func(sess *mgo.Session) error {
		//多条删除
		collection := sess.DB(m.dbname).C(collection_name)
		//err := collection.Find(query_bson).All(&res)
		info, err := collection.RemoveAll(query_bson)
		if err != nil {
			return err
		}
		//s, _ := json.Marshal(res)
		//fmt.Printf("single res is :%s\n", s)
		//fmt.Println("=========", res)
		//res = s
		s, _ := json.Marshal(info)
		fmt.Printf("change info is=====%s\n", s)
		return nil
	}
	m.Execute(query_func)
}

func (m *MZMongo) Update(collection_name string, pk string, data bson.M) {
	query_bson := bson.M{pk: data[pk]}
	query_func := func(sess *mgo.Session) error {
		//单条更新
		collection := sess.DB(m.dbname).C(collection_name)
		err := collection.Update(query_bson, data)
		if err != nil {
			return err
		}
		return nil
	}
	m.Execute(query_func)
}

func (m *MZMongo) UpdateAll(collection_name string, pk string, data bson.M) {
	query_bson := bson.M{pk: data[pk]}
	delete(data, pk)
	query_func := func(sess *mgo.Session) error {
		//单条更新
		collection := sess.DB(m.dbname).C(collection_name)
		info, err := collection.UpdateAll(query_bson, bson.M{"$set": data})
		fmt.Println(999999999999)
		if err != nil {
			panic(err)
			return err
		}
		fmt.Println(999999999999)
		s, _ := json.Marshal(info)
		fmt.Printf("UpdateAll change info is=====%s\n", s)
		return nil
	}
	m.Execute(query_func)
}

func (m *MZMongo) Upsert(collection_name string, pk string, data bson.M) {
	query_bson := bson.M{pk: data[pk]}
	query_func := func(sess *mgo.Session) error {
		//单条更新
		collection := sess.DB(m.dbname).C(collection_name)
		info, err := collection.Upsert(query_bson, data)
		if err != nil {
			return err
		}
		s, _ := json.Marshal(info)
		fmt.Printf("change info is=====%s\n", s)
		return nil
	}
	m.Execute(query_func)
}

func (m *MZMongo) DropIndex(collection_name string, keys ...string) {
	query_func := func(sess *mgo.Session) error {
		//删除索引
		collection := sess.DB(m.dbname).C(collection_name)
		err := collection.DropIndex(keys...)
		if err != nil {
			return err
		}
		return nil
	}
	m.Execute(query_func)
}

func (m *MZMongo) EnsureIndex(collection_name string, keys ...string) {
	query_func := func(sess *mgo.Session) error {
		//确保索引
		collection := sess.DB(m.dbname).C(collection_name)
		err := collection.EnsureIndexKey(keys...)
		if err != nil {
			return err
		}
		return nil
	}
	m.Execute(query_func)
}

func (m *MZMongo) GetIndexes(collection_name string) []mgo.Index {
	var indexes []mgo.Index
	query_func := func(sess *mgo.Session) error {
		//获取索引信息
		collection := sess.DB(m.dbname).C(collection_name)
		idxs, err := collection.Indexes()
		if err != nil {
			return err
		}
		indexes = idxs
		return nil
	}
	m.Execute(query_func)
	return indexes
}

func (m *MZMongo) FindAndModify(collection_name string, query bson.M, update bson.M) interface{} {
	var res interface{}
	query_func := func(sess *mgo.Session) error {
		//多条查询
		collection := sess.DB(m.dbname).C(collection_name)
		info, err := collection.Find(query).Apply(mgo.Change{
			Update:    update,
			ReturnNew: true,
		}, &res)
		if err != nil {
			panic(err)
			return err
		}
		fmt.Sprintf("MZMongo FindAndModify info========%#v", info)
		return nil
	}
	m.Execute(query_func)
	return res
}
