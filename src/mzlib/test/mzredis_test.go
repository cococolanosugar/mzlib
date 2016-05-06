package test

import (
	"encoding/json"
	"fmt"
	"gopkg.in/redis.v3"
	"mzlib/client"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func Assert(a, b interface{}, t *testing.T) {
	if a != b {
		t.Error("Assert fail!")
		fmt.Println("Assert fail!")
		fmt.Println("a=====", a, reflect.TypeOf(a).Name())
		fmt.Println("b=====", b, reflect.TypeOf(b).Name())
	} else {
		fmt.Println("Assert success!")
		fmt.Println("a=====", a, reflect.TypeOf(a).Name())
		fmt.Println("b=====", b, reflect.TypeOf(b).Name())
	}
}

func TestMZRedis(t *testing.T) {
	return

	fmt.Println("MZredis test begin.......")
	mredis := client.NewMZRedisClient("127.0.0.1", 6379, 1, "")
	key1 := "xushicai"
	value1 := "15618973660"
	// test get set delete
	mredis.Set(key1, value1, time.Minute)
	v1 := mredis.Get(key1)
	Assert(v1, value1, t)

	mredis.Delete(key1)

	v1 = mredis.Get(key1)
	Assert(v1, "", t)

	//test inrc
	key2 := "counter"
	value2 := "2"
	mredis.Delete(key2)
	mredis.Incr(key2)
	mredis.Incr(key2)
	v2 := mredis.Get(key2)
	Assert(v2, value2, t)

	//test mset mget
	key31 := "key31"
	value31 := "value31"
	key32 := "key32"
	value32 := "value32"
	key33 := "key33"
	value33 := "value33"
	pairs_string := []string{key31, value31, key32, value32, key33, value33}
	imap3 := map[string]string{
		"keyimap31": "456",
		"keyimap32": "000",
	}
	mredis.Mset(pairs_string)
	mredis.Mset(imap3)

	v3 := mredis.Mget(key31, key32, key33, "keyimap31", "keyimap32")
	Assert(v3[0], value31, t)
	Assert(v3[1], value32, t)
	Assert(v3[2], value33, t)
	Assert(v3[3], "456", t)
	Assert(v3[4], "000", t)
	//fmt.Println("v3============", v3)

	//test sorted set
	var score_set []redis.Z
	zset_name := "top"
	for i := 0; i < 10; i++ {
		tmp := redis.Z{
			Score:  float64(i),
			Member: "name" + strconv.Itoa(i),
		}
		score_set = append(score_set, tmp)
	}
	mredis.Zadd(zset_name, score_set...)

	cnt := mredis.Zcard(zset_name)
	Assert(cnt, int64(10), t)

	rangelist := mredis.Zrange(zset_name, 0, 100)
	fmt.Println("zrange======", rangelist)
	revrangelist := mredis.Zrevrange(zset_name, 0, 100)
	fmt.Println("revrangelist======", revrangelist)
	Assert(len(rangelist), 10, t)
	Assert(len(revrangelist), 10, t)

	mredis.Zincrby(zset_name, "name0", 10)
	rangelist = mredis.Zrange(zset_name, 0, 100)
	Assert(rangelist[len(rangelist)-1], "name0", t)

	score := mredis.Zscore(zset_name, "name0")
	Assert(score, float64(10), t)

	mredis.Zrem(zset_name, "name9")
	rangelist = mredis.Zrange(zset_name, 0, 100)
	Assert(len(rangelist), 9, t)
	fmt.Println("rangelistxxxxxxxx", rangelist)

	mredis.Zremrangebyrank(zset_name, 0, 3)
	rangelist = mredis.Zrange(zset_name, 0, 100)
	Assert(len(rangelist), 5, t)
	fmt.Println("rangelistxxxxxxxx", rangelist)

	mredis.Zremrangebyscore(zset_name, 6, 8)
	rangelist = mredis.Zrange(zset_name, 0, 100)
	Assert(len(rangelist), 2, t)
	fmt.Println("rangelistxxxxxxxx", rangelist)

	//test unsorted set
	set_name := "unsorted_set"
	set_arr := []string{"s1", "s2", "s3", "s4", "s5"}
	mredis.Sadd(set_name, set_arr...)

	cnt = mredis.Scard(set_name)
	Assert(cnt, int64(5), t)

	mbs := mredis.Smembers(set_name)
	fmt.Println("Smembersxxxxxxxxxxx", mbs)

	pop := mredis.Spop(set_name)
	fmt.Println("popxxxxxxxxxxx====", pop)
	pop1 := mredis.Spop(set_name)
	fmt.Println("pop1xxxxxxxxxxx====", pop1)

	mbs = mredis.Smembers(set_name)
	fmt.Println("Smembersxxxxxxxxxxx", mbs)

	mredis.Srem(set_name, "s1")

	mbs = mredis.Smembers(set_name)
	fmt.Println("Smembersxxxxxxxxxxx", mbs)

	s := mredis.PoolStats()
	js, _ := json.Marshal(s)
	fmt.Printf("PoolStatsxxxxxxxxx%s\n", js)

	fmt.Println("MZredis test end.........")

}
