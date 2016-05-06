package client

import (
	"fmt"
	"gopkg.in/redis.v3"
	//"reflect"
	"strconv"
	"strings"
	"time"
)

type MZRedis struct {
	Host     string
	Port     int
	DBIndex  int
	Password string
	Client   *redis.Client
}

func NewMZRedisClient(host string, port int, db int64, password string) *MZRedis {
	client := new(MZRedis)
	client.Init(host, port, db, password)
	return client
}

func (m *MZRedis) Init(host string, port int, db int64, password string) {
	opts := redis.Options{
		Addr:     strings.Join([]string{host, strconv.Itoa(port)}, ":"),
		DB:       db,
		Password: password,
		PoolSize: 100,
		//IdleTimeout: time.Minute * 10,
	}

	m.Client = redis.NewClient(&opts)
}

func (m *MZRedis) Get(key string) string {
	value, err := m.Client.Get(key).Result()
	if err == redis.Nil {
		return ""
	}
	checkError(err)
	return value
}

func (m *MZRedis) Set(key string, value interface{}, expiration time.Duration) {
	fmt.Printf("MZRedis Set v===============%s\n", value)
	_, err := m.Client.Set(key, value, expiration).Result()
	checkError(err)
}

func (m *MZRedis) Delete(keys ...string) {
	_, err := m.Client.Del(keys...).Result()
	checkError(err)
}

func (m *MZRedis) Incr(key string) {
	_, err := m.Client.Incr(key).Result()
	checkError(err)
}

func (m *MZRedis) Mset(pairs interface{}) {
	//fmt.Println("88888888888", reflect.TypeOf(pairs))
	switch v := pairs.(type) {
	case []string:
		tmp_pairs := pairs.([]string)
		_, err := m.Client.MSet(tmp_pairs...).Result()
		checkError(err)
		//fmt.Println("type of v====1", v)
	case map[string]string:
		imap, ok := pairs.(map[string]string)
		if !ok {
			fmt.Println("mset fail!")
			return
		}
		var tmp_pairs []string
		for k, v := range imap {
			tmp_pairs = append(tmp_pairs, k, v)
		}
		//fmt.Println("tmp_pairs======", tmp_pairs)
		//fmt.Println("type of v====2", v)
		_, err := m.Client.MSet(tmp_pairs...).Result()
		checkError(err)
	default:
		fmt.Println("unknown type of v!!!!===", v)
	}
}

func (m *MZRedis) Mget(keys ...string) []interface{} {
	arr, err := m.Client.MGet(keys...).Result()
	checkError(err)
	return arr
}

func (m *MZRedis) Zadd(set_name string, members ...redis.Z) {
	_, err := m.Client.ZAdd(set_name, members...).Result()
	checkError(err)
}

func (m *MZRedis) Zcard(set_name string) int64 {
	res, err := m.Client.ZCard(set_name).Result()
	checkError(err)
	return res
}

func (m *MZRedis) Zincrby(set_name string, member string, increment float64) {
	_, err := m.Client.ZIncrBy(set_name, increment, member).Result()
	checkError(err)
}

//start 从零开始
func (m *MZRedis) Zrange(set_name string, start, stop int64) []string {
	res, err := m.Client.ZRange(set_name, start, stop).Result()
	checkError(err)
	return res
}

//start 从零开始
func (m *MZRedis) Zrevrange(set_name string, start, stop int64) []string {
	res, err := m.Client.ZRevRange(set_name, start, stop).Result()
	checkError(err)
	return res
}

//排行id返回带分数
func (m *MZRedis) Zrangewithscores(set_name string, start, stop int64) []interface{} {
	zlist, err := m.Client.ZRangeWithScores(set_name, start, stop).Result()
	checkError(err)
	name_score_list := make([]interface{}, 0)
	for _, z := range zlist {
		m := []interface{}{z.Member, z.Score}
		name_score_list = append(name_score_list, m)
	}
	return name_score_list
}

//排行id返回带分数
func (m *MZRedis) Zrevrangewithscores(set_name string, start, stop int64) []interface{} {
	zlist, err := m.Client.ZRevRangeWithScores(set_name, start, stop).Result()
	checkError(err)
	name_score_list := make([]interface{}, 0)
	for _, z := range zlist {
		m := []interface{}{z.Member, z.Score}
		name_score_list = append(name_score_list, m)
	}
	return name_score_list
}

func (m *MZRedis) Zrank(set_name string, member string) int64 {
	res, err := m.Client.ZRank(set_name, member).Result()
	if err == redis.Nil {
		return -1
	}
	checkError(err)
	return res
}

func (m *MZRedis) Zrevrank(set_name string, member string) int64 {
	res, err := m.Client.ZRevRank(set_name, member).Result()
	if err == redis.Nil {
		return -1
	}
	checkError(err)
	return res
}

func (m *MZRedis) Zscore(set_name string, member string) float64 {
	res, err := m.Client.ZScore(set_name, member).Result()
	if err == redis.Nil {
		return -1
	}
	checkError(err)
	return res
}

func (m *MZRedis) Zrem(set_name string, members ...string) {
	_, err := m.Client.ZRem(set_name, members...).Result()
	checkError(err)
}

func (m *MZRedis) Zremrangebyrank(set_name string, start int64, stop int64) {
	_, err := m.Client.ZRemRangeByRank(set_name, start, stop).Result()
	checkError(err)
}

func (m *MZRedis) Zremrangebyscore(set_name string, min int64, max int64) {

	//_, err := m.Client.ZRemRangeByScore(set_name, strconv.Itoa(min), strconv.Itoa(max)).Result()
	_, err := m.Client.ZRemRangeByScore(set_name, strconv.FormatInt(min, 10), strconv.FormatInt(max, 10)).Result()
	checkError(err)
}

func (m *MZRedis) Sadd(set_name string, members ...string) {
	_, err := m.Client.SAdd(set_name, members...).Result()
	checkError(err)
}

func (m *MZRedis) Scard(set_name string) int64 {
	res, err := m.Client.SCard(set_name).Result()
	checkError(err)
	return res
}

func (m *MZRedis) Smembers(set_name string) []string {
	res, err := m.Client.SMembers(set_name).Result()
	checkError(err)
	return res
}

func (m *MZRedis) Spop(set_name string) string {
	res, err := m.Client.SPop(set_name).Result()
	checkError(err)
	return res
}

func (m *MZRedis) Srem(set_name string, members ...string) {
	_, err := m.Client.SRem(set_name, members...).Result()
	checkError(err)
}

func (m *MZRedis) PoolStats() *redis.PoolStats {
	s := m.Client.PoolStats()
	return s
}
