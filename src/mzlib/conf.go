package mzlib

import (
	"fmt"
)

const (
	DEV_MOD = true
)

type MySQLConfig struct {
	Host        string  `json:"host"`
	Port        int     `json:"port"`
	Uname       string  `json:"uname"`
	Pass        string  `json:"pass"`
	Vnode       uint8   `json:"vnode"`
	NodeName    string  `json:"nodename"`
	Dbname      string  `json:"dbname"`
	Charset     string  `json:"charset"`
	PoolSize    uint16  `json:"pool"`
	IdleTimeOut float64 `json:"idle"`
	MaxRetry    uint8   `json:"retry"`
}

type MongodbConfig struct {
	Host        string  `json:"host"`
	Port        int     `json:"port"`
	Uname       string  `json:"uname"`
	Pass        string  `json:"pass"`
	Vnode       uint8   `json:"vnode"`
	NodeName    string  `json:"nodename"`
	Dbname      string  `json:"dbname"`
	Charset     string  `json:"charset"`
	PoolSize    uint16  `json:"pool"`
	IdleTimeOut float64 `json:"idle"`
	MaxRetry    uint8   `json:"retry"`
}

type RedisConfig struct {
	Host        string  `json:"host"`
	Port        int     `json:"port"`
	Password    string  `json:"password"`
	Index       int64   `json:"index"`
	Vnode       uint8   `json:"vnode"`
	NodeName    string  `json:"nodename"`
	PoolSize    uint16  `json:"pool"`
	IdleTimeOut float64 `json:"idle"`
	MaxRetry    uint8   `json:"retry"`
}

var mongo_conf MongodbConfig
var redis_conf RedisConfig
var mongo_log_conf MongodbConfig
var redis_top_conf RedisConfig

func InitConfig() {

	if DEV_MOD {
		initDevConfig()
	} else {
		InitProdConfig()
	}
}

func initDevConfig() {
	fmt.Println("conf init start.........")

	mongo_conf = MongodbConfig{
		Host:     "127.0.0.1",
		Port:     27017,
		Uname:    "",
		Pass:     "",
		Dbname:   "test",
		PoolSize: 100,
	}

	redis_conf = RedisConfig{
		Host:     "127.0.0.1",
		Port:     6379,
		Index:    0,
		PoolSize: 100,
	}

	mongo_log_conf = MongodbConfig{
		Host:     "127.0.0.1",
		Port:     27017,
		Uname:    "",
		Pass:     "",
		Dbname:   "test",
		PoolSize: 100,
	}

	redis_top_conf = RedisConfig{
		Host:     "127.0.0.1",
		Port:     6379,
		Index:    0,
		PoolSize: 100,
	}

	fmt.Println("conf init end.........")
}

func InitProdConfig() {

	fmt.Println("conf init start.........")

	fmt.Println("conf init end.........")
}
