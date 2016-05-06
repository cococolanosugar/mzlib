/*
排行榜相关功能
*/

package mzlib

import (
	//"encoding/json"
	"fmt"
	//"gopkg.in/mgo.v2/bson"
	"gopkg.in/redis.v3"
	//"reflect"
	//"time"
)

//topmodel 排行榜
type TopModel struct {
	TopName string
}

func NewTopModel(top_name string) *TopModel {
	tm := new(TopModel)
	tm.Create(top_name)
	return tm
}

func (tm *TopModel) Create(top_name string) {
	//创建一个排行榜
	tm.TopName = fmt.Sprintf("top:%s", top_name)
}

func (tm *TopModel) Get(count int64, desc bool) []string {
	//取前多少名，count为数量，desc=True从大到小 desc=False从小到大,返回一个list, count<1取所有的
	var list []string
	if desc {
		list = app.top_redis.Zrevrange(tm.TopName, 0, count-1)
	} else {
		list = app.top_redis.Zrange(tm.TopName, 0, count-1)
	}
	return list
}

func (tm *TopModel) GetWithScores(count int64, desc bool) []interface{} {
	//取前多少名，count为数量，desc=True从大到小 desc=False从小到大,返回一个list, count<1取所有的
	var list []interface{}
	if desc {
		list = app.top_redis.Zrevrangewithscores(tm.TopName, 0, count-1)
	} else {
		list = app.top_redis.Zrangewithscores(tm.TopName, 0, count-1)
	}
	return list
}

func (tm *TopModel) Set(name string, score int) {
	//加一个人 正确返回1，更新返回0
	one := redis.Z{
		Score:  float64(score),
		Member: name,
	}
	app.top_redis.Zadd(tm.TopName, one)
}

func (tm *TopModel) Remove(names ...string) {
	//踢出去一部分人 返回被踢出人数 names=[name1,name2 ...]
	app.top_redis.Zrem(tm.TopName, names...)
}

func (tm *TopModel) RemoveByScore(min_score, max_score int64) {
	//踢出去一部分人 分数在min,max之间的(包括min,max) 返回个数
	app.top_redis.Zremrangebyscore(tm.TopName, min_score, max_score)
}

func (tm *TopModel) RemoveLast(count int64, desc bool) {
	//踢出最后若干个 desc=True从大到小 desc=False从小到大 返回个数
	if desc {
		cnt := tm.Count()
		app.top_redis.Zremrangebyrank(tm.TopName, cnt-count, cnt-1)
	} else {
		app.top_redis.Zremrangebyrank(tm.TopName, 0, count-1)
	}

}

func (tm *TopModel) Rank(name string, desc bool) int64 {
	//取名次 返回一个int型，desc=True从大到小 desc=False从小到大
	var rank int64
	if desc {
		rank = app.top_redis.Zrevrank(tm.TopName, name)
	} else {
		rank = app.top_redis.Zrank(tm.TopName, name)
	}
	return rank
}

func (tm *TopModel) Count() int64 {
	//返回排行榜总数
	cnt := app.top_redis.Zcard(tm.TopName)
	return cnt
}

func (tm *TopModel) Score(name string) float64 {
	//返回分数
	score := app.top_redis.Zscore(tm.TopName, name)
	return score
}
